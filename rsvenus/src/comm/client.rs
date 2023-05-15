use std::{collections::HashMap, sync::{RwLock, Arc, Mutex}, io, thread};

use super::{common::Sender, MessageChannel, ChannelEvent, Message};

#[derive(Clone)]
pub struct Client {
  id: String,
  endpoint: bool,
  senders: Arc<RwLock<HashMap<String, Box<dyn Sender>>>>,
  mailboxes: Arc<Mutex<Vec<Box<dyn Mailbox>>>>
}

impl Client {
  pub fn new(id: impl AsRef<str>) -> Self {
    Client {
      id: id.as_ref().into(),
      endpoint: true,
      senders: Arc::new(RwLock::new(HashMap::new())),
      mailboxes: Arc::new(Mutex::new(Vec::new())),
    }
  }

  pub fn new_router(id: impl AsRef<String>) -> Self {
    Client {
      id: id.as_ref().into(),
      endpoint: false,
      senders: Arc::new(RwLock::new(HashMap::new())),
      mailboxes: Arc::new(Mutex::new(Vec::new())),
    }
  }

  pub fn id(&self) -> String {
    self.id.clone()
  }

  pub fn attach(&mut self, mailbox: impl Mailbox + 'static) {
    self.mailboxes.lock().unwrap().push(Box::new(mailbox));
  }

  fn spread_message(&mut self, message: Message, allow_broadcast: bool) -> io::Result<()> {
    let mut senders = self.senders.write().unwrap();
    if let Some(receiver) = &message.receiver {
      if *receiver == self.id {
        return Ok(())
      }
      if let Some(sender) = senders.get_mut(receiver) {
        return sender.send(message).map(|_| ())
      }
    }
    if allow_broadcast {
      for (id, sender) in senders.iter_mut() {
        if *id != message.sender {
          sender.send(message.clone())?;
        }
      }
    }
    Ok(())
  }

  pub fn send(&mut self, message: Message) -> io::Result<()> {
    self.spread_message(message, true)
  }

  fn on_event(&mut self, message: Message, sender: Box<dyn Sender>) {
    {
      let mut senders = self.senders.write().unwrap();
      if !senders.contains_key(&message.sender) {
        senders.insert(message.sender.clone(), sender);
      }
    }
    self.spread_message(message, !self.endpoint).ok();
  }

  pub fn start_channel<T: Sender>(&mut self, mut channel: impl MessageChannel<T>) -> io::Result<()> {
    let receiver = channel.start()?;
    let mailboxes = self.mailboxes.clone();
    let mut client = self.clone();
    thread::spawn(move || {
      for event in receiver.iter() {
        if event.0.sender == client.id {
          continue;
        }
        client.on_event(event.0.clone(), Box::new(event.1));
        if let Some(sender) = client.senders.write().unwrap().get_mut(&event.0.sender) {
          for mailbox in mailboxes.lock().unwrap().iter_mut() {
            mailbox.notify((event.0.clone(), sender.as_mut()), client.clone())
          }
        }
      }
    });
    Ok(())
  }
}

pub trait Mailbox: Send + Sync {
  fn notify(&mut self, event: ChannelEvent<&mut dyn Sender>, client: Client);
}