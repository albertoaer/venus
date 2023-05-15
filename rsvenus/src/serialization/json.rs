use crate::comm::{MessageSerializer, Message};

#[derive(Clone, Copy)]
pub struct JsonSerializer;

impl MessageSerializer for JsonSerializer {
  fn deserialize(data: Vec<u8>) -> Result<Message, String> {
    serde_json::from_slice::<Message>(&data).map_err(|err| err.to_string())
  }

  fn serialize(message: Message) -> Result<Vec<u8>, String> {
    serde_json::to_vec(&message).map_err(|err| err.to_string())
  }
}