use std::{ops::Deref, sync::{Mutex, Arc}};

use crate::{runtime::RuntimeContext, Task};

pub struct EventContext<'a, T> {
  context: &'a RuntimeContext,
  event: &'a T
}

impl<'a, T> EventContext<'a, T> {
  pub fn event(&self) -> &T {
    &self.event
  }
}

impl<'a, T> Deref for EventContext<'a, T> {
  type Target = RuntimeContext;

  fn deref(&self) -> &Self::Target {
    &self.context
  }
}

pub type EventTask<T> = dyn FnMut(EventContext<T>) -> bool + Send + 'static;

#[derive(Clone)]
pub struct EventBuilder<T> {
  task: Option<Arc<Mutex<Box<EventTask<T>>>>>,
  event: Option<T>
}

impl<T: Send + 'static> EventBuilder<T> {
  pub fn new(task: impl Send + 'static + FnMut(EventContext<T>) -> bool) -> Self {
    EventBuilder { task: Some(Arc::new(Mutex::new(Box::new(task)))), event: None }
  }

  pub fn task(mut self, task: impl Send + 'static + FnMut(EventContext<T>) -> bool) -> Self {
    self.task = Some(Arc::new(Mutex::new(Box::new(task))));
    self
  }

  pub fn event(mut self, event: T) -> Self {
    self.event = Some(event);
    self
  }

  pub fn build(self) -> Task {
    let task = self.task.unwrap();
    let event = self.event.unwrap();
    Arc::new(Mutex::new(Box::new(move |context| task.lock().unwrap()(EventContext { context, event: &event }))))
  }
}