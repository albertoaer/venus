use std::{collections::HashMap, sync::{Arc, RwLock}};

use crate::{Runtime, comm::{self, Mailbox, Message}, EventBuilder, EventContext};

#[derive(Clone)]
pub struct MailEvent {
  pub message: Message,
}

#[derive(Clone)]
pub struct MailboxedRuntime {
  runtime: Runtime,
  responses: Arc<RwLock<HashMap<String, EventBuilder<MailEvent>>>>,
}

pub fn mailboxed<'a>(runtime: &Runtime) -> MailboxedRuntime {
  return MailboxedRuntime { runtime: runtime.clone(), responses: Arc::new(RwLock::new(HashMap::new())) }
}

impl MailboxedRuntime {
  pub fn on(&mut self, verb: impl AsRef<str>, task: impl Send + 'static + FnMut(EventContext<MailEvent>) -> bool) {
    self.responses.write().unwrap().insert(verb.as_ref().into(), EventBuilder::new(task));
  }
}

impl Mailbox for MailboxedRuntime {
  fn notify(&mut self, event: comm::ChannelEvent<&mut dyn comm::Sender>, _: comm::Client) {
    if let Some(task) = self.responses.read().unwrap().get(&event.0.verb) {
      let event = task.clone().event(MailEvent { message: event.0.clone() });
      self.runtime.launch_boxed(event.build(), None);
    }
  }
}