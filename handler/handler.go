package handler

import (
	"wsserver/msg"
	"wsserver/network"
)

func init() {
	network.Register(msg.ID_LOGIN_REQ, &msg.LoginReq{}, Login)
	network.Register(msg.ID_LOGIN_RSP, &msg.LoginRsp{}, nil)

	network.Register(msg.ID_LESSION_STATUS_UPDATE_BROADCAST, &msg.PushData{}, nil)
}
