package protocol

import "github.com/google/uuid"

var idGen = uuid.New()

type baseParticipant struct {
	id string
}

func NewParticipant() ConversationParticipant {
	return &baseParticipant{
		id: idGen.String(),
	}
}

func (*baseParticipant) GetId() ConversationId {
	panic("unimplemented")
}
