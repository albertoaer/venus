package protocol

import "encoding/json"

type jsonMessage struct {
	Sender_            ClientId          `json:"sender"`
	Receiver_          *ClientId         `json:"receiver,omitempty"`
	Timestamp_         int64             `json:"timestamp"`
	PreviousTimestamp_ *int64            `json:"previousTimestamp,omitempty"`
	Type_              MessageType       `json:"type"`
	Verb_              Verb              `json:"verb"`
	Args_              []string          `json:"args,omitempty"`
	Options_           map[string]string `json:"options,omitempty"`
	Payload_           []byte            `json:"payload,omitempty"`
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
	message := jsonMessage{
		Args_:    make([]string, 0),
		Options_: make(map[string]string, 0),
		Payload_: make([]byte, 0),
	}
	err := json.Unmarshal(packet, &message)
	return message, err
}

func (*jsonSerializer) Serialize(msg Message) ([]byte, error) {
	message := jsonMessage{
		Sender_:            msg.Sender(),
		Receiver_:          msg.Receiver(),
		Timestamp_:         msg.Timestamp(),
		PreviousTimestamp_: msg.PreviousTimestamp(),
		Type_:              msg.Type(),
		Verb_:              msg.Verb(),
		Args_:              msg.Args(),
		Options_:           msg.Options(),
		Payload_:           msg.Payload(),
	}
	return json.Marshal(message)
}
