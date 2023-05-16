use std::sync::{RwLock, Arc};

use concurrent_queue::ConcurrentQueue;

use crate::{runtime::Promise, RuntimePlan, Runtime};

pub struct SingleThreadRuntimePlan {
  queue: ConcurrentQueue<Promise>,
  on: Arc<RwLock<bool>>,
}

impl SingleThreadRuntimePlan {
  pub fn new() -> Self {
    return Self {
      queue: ConcurrentQueue::unbounded(),
      on: Arc::new(RwLock::new(false))
    }
  }
}

impl RuntimePlan for SingleThreadRuntimePlan {
  fn launch(&self, promise: Promise) {
    self.queue.push(promise).ok();
  }

  fn start(&self) {
    *self.on.write().unwrap() = true;
    while *self.on.read().unwrap() {
      if let Ok(mut task) = self.queue.pop() {
        if task.is_available() {
          task.run_once();
          if !task.is_done() {
            self.queue.push(task).ok();
          }
        } else {
          self.queue.push(task).ok();
        }
      }
    }
  }

  fn stop(&self) {
    *self.on.write().unwrap() = false;
  }
}

impl Into<Runtime> for SingleThreadRuntimePlan {
  fn into(self) -> Runtime {
    Runtime::new(self)
  }
}