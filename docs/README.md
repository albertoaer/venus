# Venus : Golang Docs

## Message

A message is the basic unit of communication. It has a very high level design intented to a fast comprenhension between both communicators.

Parameters:
- `Sender`: the client that sent the message
- `Receiver`: the client that this message is sent to, might be empty for broadcast
- `Timestamp`: the moment when the message was created
- `Args`: *string[]* the arguments required by the task
- `Options`: *string[string]* the options required by the task, mapped by name
- `Payload`: *byte[]* a generic purpose byte array

Can be created using a `MessageBuilder`

```go
func createMessage(client comm.Client) comm.Message {
  return comm NewMessageBuilder(client.GetId()).SetVerb(comm.Hi).Build()
}
```

## Client

The one that mediates in the communication process. Identified, named **ClientId** by any distributed compatible ID format, like UUID, GUID or ULID (the one used at the golang implementation).

Will manage the **channels** lifecycle and notify the **mailboxes** on any incoming message that verifies the protocol.

```go
func createClient() comm.Client {
  return comm.NewClient(utils.NewUlidIdGenerator().NextId())
}
```

## Channel

Channels are used to receive and reply messages. Can be manually handled or by a client.

```go
func mustAddTcpChannel(client comm.Client, port int) {
	tcpChannel := network.NewTcpChannel().SetPort(port).AsMessageChannel()
	if err := client.StartChannel(tcpChannel); err != nil {
		panic(err)
	}
}
```

Channels are usually closed to the protocol ruled by the client, but could not be possible to know in the first instance a tcp/udp address by a **ClientId**, so there must be a way to identify to another **Client** our existence. **OpenableChannel** is a special channel that can create a gateway to an address.

```go
func helloWorldToAddress(address net.Addr) {
  // TcpChannel is OpenableChannel[net.Addr]
  tcpChannel := network.NewTcpChannel().AsMessageChannel()
  msg := ... // generate the message
  tcpChannel.Open(address).Send(msg)
}
```

## Mailbox

The **Mailbox** interface is implemented by any type with the method `Notify(ChannelEvent, Client)` and is capable of receiving messages arrived on the server.

```go
func addSniffer(client comm.Client) {
  // The sniffer is an implemented mailbox that logs the received messages
  client.Attach(comm.NewSniffer())
}
```

## Runtime

A runtime is the task executor, the component that decides when and how to perform the tasks it is assigned to launch.

```go
func ready(c govenus.MailContext) bool {
	fmt.Println("Got hi")
	reply := comm.NewMessageBuilder(c.Event().Client.GetId()).
		SetVerb(comm.Hi).
    SetReceiver(c.Event().Message.Sender)
	c.Event().Sender.Send(reply) // Send back the reply
	return true
}

func getMailbox() comm.Mailbox {
  // Create a new single routine runtime
  runtime := govenus.NewSRRuntime()
  // Mailboxed function creates a mailbox that will match messages and tasks, launching them in the runtime
	mailbox := govenus.Mailboxed(runtime)
	mailbox.On(comm.Hi, ready)
  return mailbox // this mailbox can later be attached to a client
}
```

Single Routine Runtime performs its tasks in a single routine, as the name suggest, so requires no mutexes or concurrent techniques fo tasks common resources access.