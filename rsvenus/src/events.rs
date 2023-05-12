use std::{ops::Deref};

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

pub type EventTask<T> = dyn FnMut(EventContext<T>) -> bool;

pub struct EventBuilder<T> {
  task: Option<Box<EventTask<T>>>,
  event: Option<T>
}

impl<T: 'static> EventBuilder<T> {
  pub fn new(task: impl FnMut(EventContext<T>) -> bool + 'static) -> Self {
    EventBuilder { task: Some(Box::new(task)), event: None }
  }

  pub fn task(mut self, task: impl FnMut(EventContext<T>) -> bool + 'static) -> Self {
    self.task = Some(Box::new(task));
    self
  }

  pub fn boxed_task(mut self, task: Box<EventTask<T>>) -> Self {
    self.task = Some(task);
    self
  }

  pub fn event(mut self, event: T) -> Self {
    self.event = Some(event);
    self
  }

  pub fn build(self) -> Box<Task> {
    let mut task = self.task.unwrap();
    let event = self.event.unwrap();
    Box::new(move |context| task(EventContext { context, event: &event }))
  }
}