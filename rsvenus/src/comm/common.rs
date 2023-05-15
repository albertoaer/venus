use serde::{Serialize, Deserialize};
use std::{collections::HashMap, io, sync::mpsc::Receiver};

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Message {
  pub sender: String,
  pub receiver: Option<String>,
  pub timestamp: u128,
  pub verb: String,
  #[serde(default)]
  pub args: Vec<String>,
  #[serde(default)]
  pub options: HashMap<String, String>,
  #[serde(default)]
  pub payload: Vec<u8>
}

pub trait Sender: Send + Sync + 'static {
  fn send(&mut self, message: Message) -> io::Result<bool>;
}

pub trait MessageSerializer: Clone + Copy + Send + Sync + 'static {
  fn deserialize(data: Vec<u8>) -> Result<Message, String>;
  fn serialize(message: Message) -> Result<Vec<u8>, String>;
}

pub type ChannelEvent<T> = (Message, T);

pub trait MessageChannel<T : Sender> {
  fn start(&mut self) -> io::Result<Receiver<ChannelEvent<T>>>;
}

pub trait OpenableChannel<T : Sender, A> {
  fn open(&mut self, address: A) -> io::Result<T>;
}