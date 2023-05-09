package service

import (
	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/socketmanager"
	"time"
)

const (
	addToClusterProgressTimeout        = 30 * time.Second
	successAddToClusterProgressTimeout = 5 * time.Second
	errAddToClusterProgressTimeout     = 10 * time.Second
	removeFromClusterProgress          = 30 * time.Second
	successRemoveFromClusterProgress   = 5 * time.Second
	errRemoveFromClusterProgress       = 10 * time.Second
	metricsTimeout                     = 30 * time.Second
)

type messageWithTimeout struct {
	msg     *socketmanager.Message
	expired time.Time
}

type initMessages struct {
	metrics                   messageWithTimeout
	addToClusterProgress      map[int]messageWithTimeout
	removeFromClusterProgress map[int]messageWithTimeout
}

func newInitMessages() *initMessages {
	return &initMessages{
		metrics:                   messageWithTimeout{},
		addToClusterProgress:      make(map[int]messageWithTimeout),
		removeFromClusterProgress: make(map[int]messageWithTimeout),
	}
}

func (i *initMessages) PushMetrics(msg *socketmanager.Message) {
	i.metrics = messageWithTimeout{msg: msg, expired: time.Now().Add(metricsTimeout)}
}

func (i *initMessages) PushAddToCluster(nodeID int, msg *socketmanager.Message) {
	expired := time.Now()
	switch msg.Payload.(internal.AddNodeToClusterProgressMsg).Status {
	case internal.STATUS_ERROR:
		expired = expired.Add(errAddToClusterProgressTimeout)
	case internal.STATUS_SUCCESS:
		expired = expired.Add(successAddToClusterProgressTimeout)
	default:
		expired = expired.Add(addToClusterProgressTimeout)
	}
	i.addToClusterProgress[nodeID] = messageWithTimeout{msg: msg, expired: expired}
}

func (i *initMessages) PushRemoveFromCluster(nodeID int, msg *socketmanager.Message) {
	expired := time.Now()
	switch msg.Payload.(internal.RemoveNodeFromClusterMsg).Status {
	case internal.STATUS_ERROR:
		expired = expired.Add(errRemoveFromClusterProgress)
	case internal.STATUS_SUCCESS:
		expired = expired.Add(successRemoveFromClusterProgress)
	default:
		expired = expired.Add(removeFromClusterProgress)
	}
	i.removeFromClusterProgress[nodeID] = messageWithTimeout{msg: msg, expired: expired}
}

func (i *initMessages) GetInitMessages() []*socketmanager.Message {
	msgs := append(make([]*socketmanager.Message, 0, len(i.addToClusterProgress)+len(i.removeFromClusterProgress)+1))

	getValid := func(m map[int]messageWithTimeout) {
		for key, msg := range m {
			if time.Now().Before(msg.expired) || (!msg.msg.Sent && msg.msg.MustSent) {
				msgs = append(msgs, msg.msg)
			} else {
				delete(m, key)
			}
		}
	}

	getValid(i.addToClusterProgress)
	getValid(i.removeFromClusterProgress)

	if i.metrics.msg != nil {
		msgs = append(msgs, i.metrics.msg)
	}
	return msgs
}
