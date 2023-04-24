package protocol

import "encoding/json"

type jsonMessage struct {
	Sender_            ClientId          `json:"sender"`
	Receiver_          *ClientId         `json:"receiver,omitempty"`
	Timestamp_         int64             `json:"timestamp"`
	PreviousTimestamp_ *int64            `json:"previousTimestamp,omitempty"`
	Type_              MessageType       `json:"type"`
	Verb_              Verb              `json:"verb"`
	Args_              []string          `json:"args"`
	Options_           map[string]string `json:"options"`
	Payload_           []byte            `json:"payload"`
}

func (jm jsonMessage) Args() []string {
	return jm.Args_
}

func (jm jsonMessage) Options() map[string]string {
	return jm.Options_
}

func (jm jsonMessage) Payload() []byte {
	return jm.Payload_
}

func (jm jsonMessage) PreviousTimestamp() *int64 {
	return jm.PreviousTimestamp_
}

func (jm jsonMessage) Receiver() *ClientId {
	return jm.Receiver_
}

func (jm jsonMessage) Sender() ClientId {
	return jm.Sender_
}

func (jm jsonMessage) Timestamp() int64 {
	return jm.Timestamp_
}

func (jm jsonMessage) Type() MessageType {
	return jm.Type_
}

func (jm jsonMessage) Verb() Verb {
	return jm.Verb_
}

type jsonSerializer struct{}

func (*jsonSerializer) Deserialize(packet []byte) (Message, error) {
	message := jsonMessage{}
	err := json.Unmarshal(packet, &message)
	return message, err
}

func (*jsonSerializer) Serialize(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}
