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
	s           *websocket.Conn
	sMux        sync.Mutex
	l           *zap.Logger
	ch          <-chan msgType
	cancelFuncs []func()
	cMux        sync.Mutex
}

func NewSocketManager[msgType any](ch <-chan msgType, s *websocket.Conn, l *zap.Logger) *SocketManager[msgType] {
	sm := SocketManager[msgType]{
		l:           l,
		s:           s,
		ch:          ch,
		cancelFuncs: make([]func(), 0),
	}
	go sm.run()
	return &sm
}

func (sm *SocketManager[msgType]) SetSocket(s *websocket.Conn) {
	sm.sMux.Lock()
	defer sm.sMux.Unlock()

	if sm.s != nil {
		err := sm.s.Close()
		if err != nil {
			sm.l.Error("error closing socket", zap.String("error", err.Error()))
		}
	}

	sm.cMux.Lock()
	for _, cancelFunc := range sm.cancelFuncs {
		cancelFunc()
	}
	sm.cancelFuncs = make([]func(), 0)
	sm.cMux.Unlock()

	sm.s = s
}

func (sm *SocketManager[msgType]) run() {
	for msg := range sm.ch {
		sm.sMux.Lock()
		if sm.s == nil {
			sm.sMux.Unlock()
			continue
		}
		err := sm.s.WriteJSON(msg)
		sm.sMux.Unlock()
		if err != nil {
			if isCloseError(err) || errors.Is(err, syscall.EPIPE) {
				sm.SetSocket(nil)
				return
			}
			sm.l.Error("send to socket error", zap.String("err", err.Error()))
		}
	}
}

func (sm *SocketManager[msgType]) SendResultToSocketByTicker(period time.Duration, process func() interface{}) {
	ctx, cancel := context.WithCancel(context.Background())
	sm.cMux.Lock()
	sm.cancelFuncs = append(sm.cancelFuncs, cancel)
	sm.cMux.Unlock()

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case _ = <-ticker.C:
			sm.sMux.Lock()
			if sm.s == nil {
				sm.sMux.Unlock()
				return
			}

			msg := process()

			err := sm.s.WriteJSON(msg)
			sm.sMux.Unlock()

			if err != nil {
				if isCloseError(err) || errors.Is(err, syscall.EPIPE) {
					sm.SetSocket(nil)
					return
				}
				sm.l.Error("send to socket error", zap.String("err", err.Error()))
			}
		}

	}
}

func isCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseProtocolError, websocket.CloseUnsupportedData, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure, websocket.CloseInvalidFramePayloadData, websocket.ClosePolicyViolation, websocket.CloseMessageTooBig, websocket.CloseMandatoryExtension, websocket.CloseInternalServerErr, websocket.CloseServiceRestart, websocket.CloseTryAgainLater, websocket.CloseTLSHandshake)
}
