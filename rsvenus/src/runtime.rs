use std::rc::Rc;

pub struct RuntimeContext {
  runtime: Rc<dyn Runtime>,
  check_condition: Rc<Box<dyn Fn() -> bool>>
}

impl RuntimeContext {
  pub fn runtime(&self) -> &dyn Runtime  {
    self.runtime.as_ref()
  }

  pub fn is_available(&self) -> bool {
    (self.check_condition)()
  }
}

#[derive(Default)]
pub struct RuntimeContextBuilder {
  runtime: Option<Rc<dyn Runtime>>,
  check_condition: Option<Rc<Box<dyn Fn() -> bool>>>
}

impl RuntimeContextBuilder {
  pub fn runtime(mut self, runtime: Rc<dyn Runtime>) -> Self {
    self.runtime = Some(runtime);
    self
  }

  fn get_condition(&self) -> Rc<Box<dyn Fn() -> bool>> {
    self.check_condition.clone().unwrap_or(Rc::new(Box::new(||true)))
  }

  pub fn add_condition(mut self, condition: impl Fn() -> bool + 'static) -> Self {
    let last_cond = self.get_condition();
    let cond = Rc::new(Box::new(condition));
    self.check_condition = Some(Rc::new(Box::new(move || cond() && last_cond())));
    self
  }

  pub fn build(self) -> Option<RuntimeContext> {
    Some(RuntimeContext { runtime: self.runtime.clone()?, check_condition: self.get_condition() })
  }
}

pub type Task = dyn FnMut(&RuntimeContext) -> bool;

pub trait Promise {
  fn is_done(&self) -> bool;
  fn on_done(&self, task: Box<Task>, builder: Option<RuntimeContextBuilder>) -> Rc<dyn Promise>;
}

pub trait Runtime {
  fn new_context(&self) -> RuntimeContextBuilder;
  fn launch(&self, task: Box<Task>, builder: Option<RuntimeContextBuilder>) -> Rc<dyn Promise>;
  fn start(&self);
  fn stop(&self);
}