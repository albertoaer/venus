use std::{collections::HashMap, sync::{RwLock, Arc, Mutex}, io, thread};

use super::{common::Sender, MessageChannel, ChannelEvent, Message};

struct RegisteredSender {
  distance: u32,
  sender: Box<dyn Sender>
}

impl RegisteredSender {
  fn from_event<T: Sender>(event: ChannelEvent<T>) -> Self {
    RegisteredSender { distance: event.0.distance, sender: Box::new(event.1) }
  }

  fn get_mut(&mut self) -> &mut dyn Sender {
    self.sender.as_mut()
  }

  fn try_replace_with_event<T: Sender>(&mut self, event: ChannelEvent<T>) {
    if event.0.distance <= self.distance {
      self.distance = event.0.distance;
      self.sender = Box::new(event.1);
    }
  }
}

#[derive(Clone)]
pub struct Client {
  id: String,
  endpoint: bool,
  senders: Arc<RwLock<HashMap<String, RegisteredSender>>>,
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

  pub fn new_router(id: impl AsRef<str>) -> Self {
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
        return sender.get_mut().send(message).map(|_| ())
      }
    }
    if allow_broadcast {
      for (id, sender) in senders.iter_mut() {
        if *id != message.sender {
          sender.get_mut().send(message.clone())?;
        }
      }
    }
    Ok(())
  }

  pub fn send(&mut self, message: Message) -> io::Result<()> {
    self.spread_message(message, true)
  }

  fn on_event<T: Sender>(&mut self, event: ChannelEvent<T>) {
    let mut message = event.0.clone();
    {
      let mut senders = self.senders.write().unwrap();
      match senders.get_mut(&event.0.sender) {
        Some(sender) => sender.try_replace_with_event(event),
        None => drop(senders.insert(event.0.sender.clone(), RegisteredSender::from_event(event))),
      }
    }
    message.distance += 1;
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
        let message = event.0.clone();
        client.on_event(event);
        for mailbox in mailboxes.lock().unwrap().iter_mut() {
          mailbox.notify(message.clone(), client.clone())
        }
      }
    });
    Ok(())
  }
}

pub trait Mailbox: Send + Sync {
  fn notify(&mut self, message: Message, client: Client);
}