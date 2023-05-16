use std::{collections::HashMap, sync::{Arc, RwLock}};

use crate::{Runtime, comm::{Mailbox, Message, Client}, EventBuilder, EventContext};

#[derive(Clone)]
pub struct MailEvent {
  pub message: Message,
  pub client: Client
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
  pub fn on(&mut self, verb: impl AsRef<str>, task: impl Send + 'static + FnMut(&mut EventContext<MailEvent>) -> bool) {
    self.responses.write().unwrap().insert(verb.as_ref().into(), EventBuilder::new(task));
  }
}

impl Mailbox for MailboxedRuntime {
  fn notify(&mut self, message: Message, client: Client) {
    if let Some(task) = self.responses.read().unwrap().get(&message.verb) {
      let event = task.clone().event(MailEvent { message: message.clone(), client });
      self.runtime.launch_boxed(event.build(), None);
    }
  }
}