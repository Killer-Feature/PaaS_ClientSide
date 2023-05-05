package socketmanager

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

type SocketManager struct {
	l     *zap.Logger
	msgCh MessageChan
	ws    WS
}

func NewSocketManager(l *zap.Logger) *SocketManager {
	sm := SocketManager{
		l:     l,
		msgCh: *NewMessageChan(),
		ws:    *NewWS(l),
	}
	go sm.run()
	return &sm
}

func (sm *SocketManager) Send(msg *Message) {
	sm.msgCh.Send(msg)
}

func (sm *SocketManager) HasWS() bool {
	return sm.ws.hasWS
}

func (sm *SocketManager) AddWS(newWS *websocket.Conn, initMessages []*Message) {
	err := sm.ws.Init(newWS, initMessages)
	if err != nil {
		return
	}
	sm.ws.Add(newWS)
}

func (sm *SocketManager) run() {
	defer func() {
		err := recover()
		if err != nil {
			sm.l.Error("panic recovered", zap.Any("panic with", err))
		}
	}()
	for msg := range sm.msgCh.Chan() {
		sm.ws.WriteJSON(msg)
		msg.Sent = true
	}
}

func (sm *SocketManager) RunByTicker(period time.Duration, process func() *Message) {
	process = sm.recoveryMW(process)
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sm.msgCh.Send(process())
				if !sm.ws.hasWS {
					return
				}
			}
		}
	}()
}

func (sm *SocketManager) recoveryMW(next func() *Message) func() *Message {
	return func() *Message {
		defer func() {
			err := recover()
			if err != nil {
				sm.l.Error("panic recovered", zap.Any("panic with", err))
			}
		}()
		return next()
	}
}
