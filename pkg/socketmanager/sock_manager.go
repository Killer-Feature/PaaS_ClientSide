package socketmanager

import (
	"go.uber.org/zap"
	"sync"
)

type Socket interface {
	WriteJSON(interface{}) error
	Close() error
}

type SocketManager[msgType any] struct {
	s  Socket
	mu sync.Mutex
	l  *zap.Logger
	ch <-chan msgType
}

func NewSocketManager[msgType any](ch <-chan msgType, s Socket, l *zap.Logger) *SocketManager[msgType] {
	sm := SocketManager[msgType]{
		l:  l,
		s:  s,
		ch: ch,
	}
	go sm.run()
	return &sm
}

func (sm *SocketManager[msgType]) SetSocket(s Socket) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.s != nil {
		err := sm.s.Close()
		if err != nil {
			sm.l.Error("error closing socket", zap.String("error", err.Error()))
		}
	}
	sm.s = s
}

func (sm *SocketManager[msgType]) run() {
	for msg := range sm.ch {
		sm.mu.Lock()
		if sm.s == nil {
			sm.mu.Unlock()
			continue
		}
		err := sm.s.WriteJSON(msg)
		sm.mu.Unlock()
		if err != nil {
			sm.l.Error("send to socket error", zap.String("err", err.Error()))
		}
	}
}
