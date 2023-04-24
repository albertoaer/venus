package protocol

const (
	Ping Verb = "PING"
	Find Verb = "FIND"
	Id   Verb = "ID"
)

const (
	MESSAGE_TYPE_PERFORM   MessageType = 1
	MESSAGE_TYPE_INFO      MessageType = 2
	MESSAGE_TYPE_BROADCAST MessageType = 3
)
