package network

import (
	"net"
	"reflect"
	"wsserver/log"
)

type AgentInf interface {
	Run()
	OnClose()
	WriteMsg(msg interface{})
	RemoteAddr() net.Addr
}

type Agent struct {
	Conn     Conn
	Svr      *WsServer
	UserData interface{}
}

func (a *Agent) Run() {

	for {

		data, err := a.Conn.ReadMsg()

		if err != nil {
			log.Debug("read message out:", err)
			log.Debug("read message addr:", a.Conn.RemoteAddr())
			break
		}

		log.Debug("read msg:", string(data))

		msg_id, msg, err := Default_Processor.Unmarshal(data)

		if err != nil {
			log.Error("read msg Unmarshal error:", err)
			continue
		}

		err = Default_Processor.Route(msg_id, msg, a)

		if err != nil {
			log.Error("read msg Route error:", err)
			continue
		}
	}

}

func (a *Agent) WriteMsg(msg interface{}) {
	if Default_Processor != nil {

		data, err := Default_Processor.Marshal(msg)
		if err != nil {
			log.Error("Agent marshal message:", reflect.TypeOf(msg), "error:", err)
			return
		}

		log.Debug("Agent WriteMsg data:", string(data))

		err = a.Conn.WriteMsg(data)
		if err != nil {
			log.Error("Agent write message:", reflect.TypeOf(msg), "error:", err)
		}

		log.Debug("Agent WriteMsg msg :", msg)
	}
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.Conn.RemoteAddr()
}

func (a *Agent) OnClose() {

}

func (a *Agent) Close() {
	a.Conn.Close()
}

func (a *Agent) Destroy() {
	a.Conn.Destroy()
}
