use super::{Mailbox, Message, Client};

pub enum Sniffer {
  Default,
  Custom(Box<dyn Fn(Message) + Send + Sync + 'static>)
}

impl Sniffer {
  pub fn custom<T : Fn(Message) + Send + Sync + 'static>(action: T) -> Self {
    Self::Custom(Box::new(action))
  }
}

impl Default for Sniffer {
  fn default() -> Self {
    Self::Default
  }
}

impl Mailbox for Sniffer {
  fn notify(&mut self, message: Message, _: Client) {
    match self {
      Self::Default => println!("{:?}", message),
      Self::Custom(action) => action(message),
    }
  }
}