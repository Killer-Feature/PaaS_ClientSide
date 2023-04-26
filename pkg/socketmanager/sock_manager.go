package socketmanager

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"syscall"
	"time"
)

type SocketManager[msgType any] struct {
	s                       []*websocket.Conn
	sMux                    sync.Mutex
	fMux                    sync.Mutex
	l                       *zap.Logger
	ch                      <-chan msgType
	cMux                    sync.Mutex
	cancelFunc              func()
	forceSendResultToSocket func()
}

func (sm *SocketManager[msgType]) ForceSendResultToSocket() {
	sm.fMux.Lock()
	defer sm.fMux.Unlock()
	if sm.forceSendResultToSocket != nil {
		sm.forceSendResultToSocket()
	}
}

func NewSocketManager[msgType any](ch <-chan msgType, l *zap.Logger) *SocketManager[msgType] {
	sm := SocketManager[msgType]{
		l:  l,
		s:  make([]*websocket.Conn, 0),
		ch: ch,
	}
	go sm.run()
	return &sm
}

func (sm *SocketManager[msgType]) AddSocket(s *websocket.Conn) {
	sm.sMux.Lock()
	defer sm.sMux.Unlock()

	sm.s = append(sm.s, s)
}

func (sm *SocketManager[msgType]) writeJSON(msg interface{}) {
	sm.sMux.Lock()
	defer sm.sMux.Unlock()
	for i, sock := range sm.s {
		err := sock.WriteJSON(msg)
		if err != nil {
			if isCloseError(err) || errors.Is(err, syscall.EPIPE) {
				sm.s = append(sm.s[:i], sm.s[i+1:]...)
			}
			if len(sm.s) == 0 {
				sm.cMux.Lock()
				if sm.cancelFunc != nil {
					sm.cancelFunc()
					sm.cancelFunc = nil
				}
				sm.cMux.Unlock()
				return
			}
			sm.l.Error("send to socket error", zap.String("err", err.Error()))
		}
	}
}

func (sm *SocketManager[msgType]) run() {
	for msg := range sm.ch {
		sm.writeJSON(msg)
	}
}

func (sm *SocketManager[msgType]) CountSockets() int {
	sm.sMux.Lock()
	defer sm.sMux.Unlock()
	return len(sm.s)
}

func (sm *SocketManager[msgType]) SendResultToSocketByTicker(period time.Duration, process func() interface{}) {
	ctx, cancel := context.WithCancel(context.Background())
	sm.cMux.Lock()
	sm.cancelFunc = cancel
	sm.cMux.Unlock()

	doProcessAndSend := func() {

		var msg interface{}
		for i := 0; i < 3; i++ {
			msg = process()
			if msg != nil {
				break
			}
		}
		if msg == nil {
			return
		}

		sm.writeJSON(msg)
	}

	forceChan := make(chan struct{})

	sm.fMux.Lock()
	sm.forceSendResultToSocket = func() {
		forceChan <- struct{}{}
	}
	sm.fMux.Unlock()

	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-forceChan:
				if sm.CountSockets() == 0 {
					return
				}
				doProcessAndSend()
			case <-ticker.C:
				if sm.CountSockets() == 0 {
					return
				}
				doProcessAndSend()
			}
		}
	}()
}

func isCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseProtocolError, websocket.CloseUnsupportedData, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure, websocket.CloseInvalidFramePayloadData, websocket.ClosePolicyViolation, websocket.CloseMessageTooBig, websocket.CloseMandatoryExtension, websocket.CloseInternalServerErr, websocket.CloseServiceRestart, websocket.CloseTryAgainLater, websocket.CloseTLSHandshake)
}
