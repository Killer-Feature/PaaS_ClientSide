package socketmanager

import (
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"syscall"
)

type WS struct {
	s     []*websocket.Conn
	mu    sync.Mutex
	l     *zap.Logger
	hasWS bool
}

func NewWS(l *zap.Logger) *WS {
	return &WS{
		s:     make([]*websocket.Conn, 0),
		l:     l,
		hasWS: false,
	}
}

func (ws *WS) Add(newWs *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.s = append(ws.s, newWs)
	ws.hasWS = true
}

func (ws *WS) Init(sock *websocket.Conn, initMessages []*Message) error {
	for _, msg := range initMessages {
		err := sock.WriteJSON(*msg)
		if err != nil {
			if !isCloseError(err) && !errors.Is(err, syscall.EPIPE) {
				ws.l.Error("send to socket error", zap.String("err", err.Error()))
			} else {
				return err
			}
		}
	}
	return nil
}

func (ws *WS) WriteJSON(msg *Message) {
	if !ws.hasWS {
		return
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()

	openedSockConn := make([]*websocket.Conn, 0, len(ws.s))

	for _, sock := range ws.s {
		err := sock.WriteJSON(msg)
		if err != nil && !isCloseError(err) && !errors.Is(err, syscall.EPIPE) {
			ws.l.Error("send to socket error", zap.String("err", err.Error()))
		}
		if err == nil || (!isCloseError(err) && !errors.Is(err, syscall.EPIPE)) {
			openedSockConn = append(openedSockConn, sock)
		}
	}
	ws.s = openedSockConn
	if len(ws.s) == 0 {
		ws.hasWS = false
	}
}

func isCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseProtocolError, websocket.CloseUnsupportedData, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure, websocket.CloseInvalidFramePayloadData, websocket.ClosePolicyViolation, websocket.CloseMessageTooBig, websocket.CloseMandatoryExtension, websocket.CloseInternalServerErr, websocket.CloseServiceRestart, websocket.CloseTryAgainLater, websocket.CloseTLSHandshake)
}
