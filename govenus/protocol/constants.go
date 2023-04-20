package protocol

const (
	Ping Verb = "PING"
	Find Verb = "FIND"
	Id   Verb = "ID"
)

const (
	MESSAGE_RESOLUTION_ORDERED   MessageResolutionMethod = 0
	MESSAGE_RESOLUTION_UNORDERED MessageResolutionMethod = 1
)

const (
	MESSAGE_TYPE_BEGIN   MessageType = 0
	MESSAGE_TYPE_PERFORM MessageType = 1
	MESSAGE_TYPE_INFO    MessageType = 2
)
