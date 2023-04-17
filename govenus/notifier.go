package govenus

import (
	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

type Notifier struct {
	mailboxes   *utils.ConcurrentArray[MailBox]
	provider    protocol.PackageProvider
	participant protocol.ConversationParticipant
	serializer  protocol.MessageSerializer
}

func NewNotifier(provider protocol.PackageProvider) *Notifier {
	participant := protocol.NewParticipant()
	return &Notifier{
		mailboxes:   utils.NewArray[MailBox](20),
		provider:    provider,
		participant: participant,
		serializer:  protocol.NewMessageSerializerV1(participant),
	}
}

func (notifier *Notifier) Attach(mailbox MailBox) *Notifier {
	notifier.mailboxes.Add(mailbox)
	return notifier
}

func (notifier *Notifier) Detach(mailbox MailBox) *Notifier {
	notifier.mailboxes.Remove(mailbox)
	return notifier
}

func (notifier *Notifier) Start() (err error) {
	err = notifier.provider.Start()
	if err == nil {
		for {
			emitter := notifier.provider.Emitter()
			packet := <-emitter
			message := notifier.serializer.Deserialize(packet.Data)
			notifier.mailboxes.ForEach(func(_ int, mb MailBox) {
				mb.Notify(message)
			})
		}
	}
	return
}
