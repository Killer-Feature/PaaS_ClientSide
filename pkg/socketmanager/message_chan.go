package socketmanager

import (
	"sync"
)

type MessageType string

type Message struct {
	Payload  interface{} `json:"payload"`
	Type     MessageType `json:"type"`
	Sent     bool        `json:"-"`
	MustSent bool        `json:"-"`
}

type MessageChan struct {
	ch chan *Message
	mu sync.Mutex
}

func NewMessageChan() *MessageChan {
	return &MessageChan{ch: make(chan *Message)}
}

func (mc *MessageChan) Send(msg *Message) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if msg == nil {
		return
	}
	mc.ch <- msg
}

func (mc *MessageChan) Chan() <-chan *Message {
	return mc.ch
}
