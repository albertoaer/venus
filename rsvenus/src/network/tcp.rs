use std::{
  net::{TcpStream, TcpListener, SocketAddr},
  io::{self, Write, Read},
  sync::{mpsc, RwLock, Arc, Mutex},
  collections::HashMap,
  marker::PhantomData,
  thread
};

use crate::comm::{MessageSerializer, Sender, ChannelEvent, Message, OpenableChannel, MessageChannel};

pub struct TcpSender<M : MessageSerializer>(TcpStream, PhantomData<M>);

impl<M : MessageSerializer> TcpSender<M> {
  pub fn try_clone(&self) -> io::Result<Self> {
    self.0.try_clone().map(|stream| TcpSender(stream, PhantomData))
  }
}

impl<M : MessageSerializer> Sender for TcpSender<M> {
  fn send(&mut self, message: Message) -> io::Result<bool> {
    let bin = M::serialize(message).map_err(|err| io::Error::new(io::ErrorKind::Other, err))?;
    self.0.write(&(bin.len() as u64).to_le_bytes())?;
    self.0.write(&bin)?;
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
    stream: TcpStream,
  ) -> io::Result<TcpSender<M>> {
    let tcp_sender = TcpSender(stream, PhantomData);
    self.connections.write().unwrap().insert(address, tcp_sender.try_clone()?);
    let mut shared_sender = tcp_sender.try_clone()?;
    let shared_conns = self.connections.clone();
    let channel = self.channel.clone();
    thread::spawn(move || {
      loop {
        let mut size_buffer = [0; (u64::BITS >> 3) as usize];
        let size = match shared_sender.0.read(&mut size_buffer) {
          Ok(_) => u64::from_le_bytes(size_buffer),
          Err(_) => break,
        };
        let mut message_buffer = vec![0; size as usize];
        if let Err(_) = shared_sender.0.read(&mut message_buffer) {
          break
        }
        let message = match M::deserialize(message_buffer) {
          Ok(message) => message,
          Err(_) => continue,
        };
        if let Ok(cloned_sender) = shared_sender.try_clone() {
          channel.lock().unwrap().as_mut().expect("channel not initialized").send((message, cloned_sender)).ok();
        }
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

    let connections = self.connections.clone();

    thread::spawn(move || {
      for incoming in listener.incoming() {
        if let Ok(stream) = incoming {
          if let Ok(mut conns) = connections.write() {
            let address = stream.peer_addr().unwrap();
            let sender = TcpSender(stream, PhantomData);
            conns.insert(address, sender);
          }
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
        return sender.try_clone()
      }
    }
    self.handle_tcp(address, TcpStream::connect(address)?)
  }
}