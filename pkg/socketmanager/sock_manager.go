package socketmanager

import (
	"go.uber.org/zap"
)

type Socket interface {
	WriteJSON(interface{}) error
	Close() error
}

type SocketManager[msgType any] struct {
	s  Socket
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
		if sm.s == nil {
			continue
		}
		err := sm.s.WriteJSON(msg)
		if err != nil {
			sm.l.Error("send to socket error", zap.String("err", err.Error()))
		}
	}
}
