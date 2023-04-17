package protocol

type messageSerializerV1 struct {
	participant ConversationParticipant
}

func NewMessageSerializerV1(participant ConversationParticipant) MessageSerializer {
	return &messageSerializerV1{
		participant: participant,
	}
}

func (*messageSerializerV1) Deserialize([]byte) Message {
	panic("unimplemented")
}

func (*messageSerializerV1) Serialize(msg Message) []byte {
	panic("unimplemented")
}
