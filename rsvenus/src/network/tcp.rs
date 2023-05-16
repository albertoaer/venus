use std::{
  net::{TcpStream, TcpListener, SocketAddr},
  io::{self, Write, Read},
  sync::{mpsc, RwLock, Arc, Mutex},
  collections::HashMap,
  marker::PhantomData,
  thread
};

use crate::comm::{MessageSerializer, Sender, ChannelEvent, Message, OpenableChannel, MessageChannel};

#[derive(Clone)]
pub struct TcpSender<M : MessageSerializer>(Arc<Mutex<TcpStream>>, PhantomData<M>);

impl<M : MessageSerializer> TcpSender<M> {
  pub fn new(stream: TcpStream) -> Self {
    TcpSender(Arc::new(Mutex::new(stream)), PhantomData)
  }
}

impl<M : MessageSerializer> Sender for TcpSender<M> {
  fn send(&mut self, message: &Message) -> io::Result<bool> {
    let mut socket = self.0.lock().unwrap();
    let bin = M::serialize(message).map_err(|err| io::Error::new(io::ErrorKind::Other, err))?;
    socket.write(&(bin.len() as u64).to_le_bytes())?;
    socket.write(&bin)?;
    Ok(false)
  }
}

pub const DEFAULT_PORT: i32 = 27103;

#[derive(Clone)]
pub struct TcpChannel<M : MessageSerializer> {
  address: String,
  connections: Arc<RwLock<HashMap<SocketAddr, TcpSender<M>>>>,
  channel: Arc<Mutex<Option<mpsc::Sender<ChannelEvent<TcpSender<M>>>>>>,
  serializer: PhantomData<M>,
}

impl<M : MessageSerializer> TcpChannel<M> {
  pub fn new(port: Option<i32>) -> TcpChannel<M> {
    let port = port.unwrap_or(DEFAULT_PORT);
    TcpChannel {
      address: format!("0.0.0.0:{port}"),
      connections: Arc::new(RwLock::new(HashMap::new())),
      channel: Arc::new(Mutex::new(None)),
      serializer: PhantomData,
    }
  }

  fn handle_tcp(
    &mut self,
    address: SocketAddr,
    mut stream: TcpStream,
  ) -> io::Result<TcpSender<M>> {
    let tcp_sender = TcpSender::new(stream.try_clone()?);
    self.connections.write().unwrap().insert(address, tcp_sender.clone());
    let shared_conns = self.connections.clone();
    let shared_sender = tcp_sender.clone();
    let channel = self.channel.clone();
    thread::spawn(move || {
      loop {
        let mut size_buffer = [0; (u64::BITS >> 3) as usize];
        let size = match stream.read(&mut size_buffer) {
          Ok(_) => u64::from_le_bytes(size_buffer),
          Err(_) => break,
        };
        let mut message_buffer = vec![0; size as usize];
        if let Err(_) = stream.read(&mut message_buffer) {
          break
        }
        let message = match M::deserialize(message_buffer) {
          Ok(message) => message,
          Err(_) => continue,
        };
        channel.lock().unwrap().as_mut().expect("channel not initialized")
          .send((message, shared_sender.clone())).ok();
      }
      shared_conns.write().and_then(|mut conns| Ok(conns.remove(&address))).ok()
    });
    Ok(tcp_sender)
  }
}

impl<M : MessageSerializer> MessageChannel<TcpSender<M>> for TcpChannel<M> {
  fn start(&mut self) -> io::Result<mpsc::Receiver<ChannelEvent<TcpSender<M>>>> {
    let listener = TcpListener::bind(self.address.clone())?;

    let (sender, receiver) = mpsc::channel();

    *self.channel.lock().unwrap() = Some(sender);

    let mut shared = self.clone();

    thread::spawn(move || {
      for incoming in listener.incoming() {
        if let Ok(stream) = incoming {
          let address = stream.peer_addr().unwrap();
          shared.handle_tcp(address, stream).ok();
        }
      }
    });

    Ok(receiver)
  }
}

impl<M : MessageSerializer> OpenableChannel<TcpSender<M>, SocketAddr> for TcpChannel<M> {
  fn open(&mut self, address: SocketAddr) -> io::Result<TcpSender<M>> {
    {
      let conns = self.connections.read()
        .map_err(|err| io::Error::new(io::ErrorKind::Other, err.to_string()))?;
      if let Some(sender) = conns.get(&address) {
        return Ok(sender.clone())
      }
    }
    self.handle_tcp(address, TcpStream::connect(address)?)
  }
}