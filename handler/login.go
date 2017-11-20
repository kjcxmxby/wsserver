package handler

import (
	"wsserver/log"
	"wsserver/msg"
	"wsserver/network"
	"wsserver/user"
)

type LoginReq struct {
	Sid uint32 `json:"sid"`
}

func Login(args []interface{}) {

	req := args[0].(*msg.LoginReq)
	agent := args[1].(*network.Agent)

	user := agent.UserData.(*user.User)

	log.Debug("Login user req:", req)

	user.Uid = req.Sid

	var rsp msg.LoginRsp

	rsp.Status = 0
	rsp.Err = ""

	log.Debug("Login Process end")

	agent.WriteMsg(&rsp)
}
