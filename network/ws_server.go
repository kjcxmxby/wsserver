package network

import (
	"net"
	"net/http"
	"sync"
	"time"
	"wsserver/log"
	"wsserver/router"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	maxConnNum      int
	pendingWriteNum int
	maxMsgLen       uint32
	newAgent        func(*WsConn) AgentInf
	closeAgent      func(AgentInf)
	upgrader        websocket.Upgrader
	conns           WsConnSet
	mutexConns      sync.Mutex
	wg              sync.WaitGroup
}

func (handler *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Debug("req url path :", r.URL.Path)

	r.ParseForm()

	log.Debug("xxxx", r.Form)

	if r.URL.Path != "/" && r.URL.Path != "/socket.io/" && r.Method != "GET" {
		router.RouterReq(r.URL.Path, r, w)
	} else {

		conn, err := handler.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debug("upgrade error: %v", err)
			return
		}

		conn.SetReadLimit(int64(handler.maxMsgLen))

		handler.wg.Add(1)
		defer handler.wg.Done()

		handler.mutexConns.Lock()
		if handler.conns == nil {
			handler.mutexConns.Unlock()
			conn.Close()
			return
		}
		if len(handler.conns) >= handler.maxConnNum {
			handler.mutexConns.Unlock()
			conn.Close()
			log.Debug("too many connections")
			return
		}
		handler.conns[conn] = struct{}{}
		handler.mutexConns.Unlock()

		wsConn := newWSConn(conn, handler.pendingWriteNum, handler.maxMsgLen)
		agent := handler.newAgent(wsConn)
		agent.Run()

		// cleanup
		handler.closeAgent(agent)
		wsConn.Close()
		handler.mutexConns.Lock()
		delete(handler.conns, conn)
		handler.mutexConns.Unlock()

		agent.OnClose()
	}
}

type WsServer struct {
	Addr            string
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	HTTPTimeout     time.Duration
	NewAgent        func(*WsConn) AgentInf
	CloseAgent      func(AgentInf)
	Ln              net.Listener
	handler         *WsHandler
}

func (server *WsServer) Start() {

	log.Debug("WsServer begin start!!!!")

	lis, err := net.Listen("tcp", server.Addr)

	if err != nil {
		log.Error("WsServer listen eror :", err, "addr :", server.Addr)
		log.Panic("WsServer listen eror :", err, "addr :", server.Addr)
		return
	}

	log.Debug("WsServer listen end")

	server.Ln = lis

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Info("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		log.Info("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 4096
		log.Info("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		log.Info("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}
	if server.NewAgent == nil {
		log.Error("NewAgent must not be nil")
		return
	}

	server.handler = &WsHandler{
		maxConnNum:      server.MaxConnNum,
		pendingWriteNum: server.PendingWriteNum,
		maxMsgLen:       server.MaxMsgLen,
		newAgent:        server.NewAgent,
		closeAgent:      server.CloseAgent,
		conns:           make(WsConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(server.Ln)

	log.Info("WsServer Start succ !!!! listent :", server.Addr)
}

func (server *WsServer) Close() {

	server.Ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
