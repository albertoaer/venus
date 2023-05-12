use std::{rc::Rc, cell::RefCell, sync::{RwLock, Arc, Mutex}};

use concurrent_queue::ConcurrentQueue;

use crate::runtime::{Task, RuntimeContext, Promise, Runtime, RuntimeContextBuilder};

struct FuncPromise {
  task: RefCell<Box<Task>>,
  done: Arc<Mutex<bool>>,
  context: RuntimeContext
}

impl FuncPromise {
  fn new(task: Box<Task>, context: RuntimeContext) -> Self {
    Self { task: RefCell::new(task), done: Arc::new(false.into()), context }
  }
  
  fn run_once(&self) {
    let mut done = self.done.lock().unwrap();
    if !*done {
      *done = self.task.borrow_mut()(&self.context);
    }
  }
}

impl Promise for FuncPromise {
  fn is_done(&self) -> bool {
    *self.done.lock().unwrap()
  }

  fn on_done(&self, task: Box<Task>, builder: Option<RuntimeContextBuilder>) -> Rc<dyn Promise> {
    let done = self.done.clone();
    let builder = builder.unwrap_or(self.context.runtime().new_context()).add_condition(move || *done.lock().unwrap());
    self.context.runtime().launch(task, Some(builder))
  }
}

struct SingleThreadRuntimeState {
  queue: ConcurrentQueue<Rc<FuncPromise>>,
  on: RwLock<bool>,
}

impl SingleThreadRuntimeState {
  fn new() -> Self {
    return Self {
      queue: ConcurrentQueue::unbounded(),
      on: RwLock::new(false)
    }
  }
}

#[derive(Clone)]
pub struct SingleThreadRuntime {
  state: Rc<SingleThreadRuntimeState>
}

impl SingleThreadRuntime {
  pub fn new() -> Self {
    Self { state: Rc::new(SingleThreadRuntimeState::new()) }
  }
}

impl Runtime for SingleThreadRuntime {
  fn new_context(&self) -> RuntimeContextBuilder {
    RuntimeContextBuilder::default().runtime(Rc::new(self.clone()))
  }

  fn launch(&self, task: Box<Task>, builder: Option<RuntimeContextBuilder>) -> Rc<dyn Promise> {
    let context = builder.map(|b| b.runtime(Rc::new(self.clone())))
      .unwrap_or(self.new_context()).build().unwrap();
    let promise = Rc::new(FuncPromise::new(task, context));
    self.state.queue.push(promise.clone()).ok();
    promise
  }

  fn start(&self) {
    *self.state.on.write().unwrap() = true;
    while *self.state.on.read().unwrap() {
      if let Ok(task) = self.state.queue.pop() {
        if task.context.is_available() {
          task.run_once();
          if !task.is_done() {
            self.state.queue.push(task).ok();
          }
        } else {
          self.state.queue.push(task).ok();
        }
      }
    }
  }

  fn stop(&self) {
    *self.state.on.write().unwrap() = false;
  }
}