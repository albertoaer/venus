use std::sync::{Arc, Mutex};

#[derive(Clone)]
pub struct RuntimeContext {
  runtime: Runtime,
  check_conditions: Arc<Mutex<Vec<Box<dyn Fn() -> bool + Send>>>>
}

impl RuntimeContext {
  pub fn runtime(&mut self) -> &mut Runtime  {
    &mut self.runtime
  }

  pub fn is_available(&self) -> bool {
    self.check_conditions.lock().unwrap().iter().all(|x| x())
  }
}

#[derive(Default)]
pub struct RuntimeContextBuilder {
  runtime: Option<Runtime>,
  check_conditions: Vec<Box<dyn Fn() -> bool + Send>>
}

impl RuntimeContextBuilder {
  pub fn runtime(mut self, runtime: Runtime) -> Self {
    self.runtime = Some(runtime);
    self
  }

  pub fn add_condition(mut self, condition: impl Fn() -> bool + Send + 'static) -> Self {
    self.check_conditions.push(Box::new(condition));
    self
  }

  pub fn build(self) -> Option<RuntimeContext> {
    Some(RuntimeContext {
      runtime: self.runtime.clone()?,
      check_conditions: Arc::new(Mutex::new(self.check_conditions)),
    })
  }
}

pub type Task = Arc<Mutex<Box<dyn Send + 'static + FnMut(&mut RuntimeContext) -> bool>>>;

#[derive(Clone)]
pub struct Promise {
  task: Task,
  done: Arc<Mutex<bool>>,
  context: RuntimeContext,
}

impl Promise {
  pub(super) fn run_once(&mut self) {
    let mut done = self.done.lock().unwrap();
    if !*done {
      *done = self.task.lock().unwrap()(&mut self.context)
    }
  }

  pub(super) fn is_available(&self) -> bool {
    self.context.is_available()
  }

  pub fn is_done(&self) -> bool {
    *self.done.lock().unwrap()
  }

  pub fn on_done(
    &mut self,
    task: impl FnMut(&mut RuntimeContext) -> bool + 'static + Send,
    builder: Option<RuntimeContextBuilder>,
  ) -> Promise {
    self.context.runtime.launch(task, builder)
  }
}

pub trait RuntimePlan: Send + Sync {
  fn launch(&self, promise: Promise);
  fn start(&self);
  fn stop(&self);
}

#[derive(Clone)]
pub struct Runtime {
  plan: Arc<Box<dyn RuntimePlan>>,
}

impl Runtime {
  pub fn new(plan: impl RuntimePlan + 'static) -> Self {
    Runtime { plan: Arc::new(Box::new(plan)) }
  }

  pub fn new_context(&self) -> RuntimeContextBuilder {
    RuntimeContextBuilder::default().runtime(self.clone())
  }

  pub fn launch(
    &mut self,
    task: impl FnMut(&mut RuntimeContext) -> bool + 'static + Send,
    builder: Option<RuntimeContextBuilder>,
  ) -> Promise {
    let context = builder.map(|b| b.runtime(self.clone()))
      .unwrap_or(self.new_context()).build().unwrap();
    let promise = Promise {
      context,
      done: Arc::new(Mutex::new(false)),
      task: Arc::new(Mutex::new(Box::new(task))),
    };
    self.plan.launch(promise.clone());
    promise
  }

  pub fn launch_boxed(
    &mut self,
    task: Task,
    builder: Option<RuntimeContextBuilder>,
  ) -> Promise {
    let context = builder.map(|b| b.runtime(self.clone()))
      .unwrap_or(self.new_context()).build().unwrap();
    let promise = Promise {
      context,
      done: Arc::new(Mutex::new(false)),
      task,
    };
    self.plan.launch(promise.clone());
    promise
  }

  pub fn start(&mut self) {
    self.plan.start()
  }
  
  pub fn stop(&mut self) {
    self.plan.stop()
  }
}