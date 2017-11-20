package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"wsserver/comm"
	"wsserver/conf"
	"wsserver/controller"
	_ "wsserver/handler"
	"wsserver/log"
	"wsserver/network"
	"wsserver/router"
	"wsserver/user"
)

func init() {
	router.Register("/socketapi/v4/lesson/status/update", controller.LessonStatusUpdate)
}

func main() {

	//closeSig := make(chan bool)

	var wsServer *network.WsServer

	wsServer = new(network.WsServer)
	wsServer.Addr = conf.SvrConf.WSAddr
	wsServer.MaxConnNum = conf.SvrConf.MaxConnNum
	wsServer.PendingWriteNum = conf.SvrConf.PendingWriteNum
	wsServer.MaxMsgLen = conf.SvrConf.MaxMsgLen
	wsServer.HTTPTimeout = time.Duration(conf.SvrConf.HTTPTimeout * int(time.Second))

	wsServer.NewAgent = func(con *network.WsConn) network.AgentInf {
		a := &network.Agent{Conn: con}
		u := user.NewUser(a)
		a.UserData = u
		user.InsertUser(u)
		return a
	}
	wsServer.CloseAgent = func(agent network.AgentInf) {
		user.RemoveUser(agent)
	}

	log.Debug("//////////version:", comm.Version, "//////////")

	if wsServer != nil {
		wsServer.Start()
	}

	//<-closeSig

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	signal.Ignore(syscall.SIGPIPE)
	sig := <-c

	log.Info("wsServer closing down signal:", sig, "!!!!!")

	if wsServer != nil {
		wsServer.Close()
	}

	log.Info("wsServer closing down !!!!!")
	log.Info("///////////////////////////////////")
}
