package protocol

import "encoding/json"

type jsonSerializer struct{}

func (*jsonSerializer) Deserialize(packet []byte) (message Message, err error) {
	message = Message{}
	err = json.Unmarshal(packet, &message)
	return
}

func (*jsonSerializer) Serialize(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}
