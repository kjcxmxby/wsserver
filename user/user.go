package user

import (
	"sync"
	"wsserver/log"
	"wsserver/network"
)

func init() {
	Users = make(UserSet)
}

type User struct {
	Uid   string
	Agent network.AgentInf
}

type UserSet map[network.AgentInf]*User

var (
	Users UserSet
	mu    sync.Mutex
)

func NewUser(a network.AgentInf) *User {
	u := &User{Agent: a}
	return u
}

func InsertUser(u *User) {
	mu.Lock()
	Users[u.Agent] = u
	mu.Unlock()
	log.Debug("new client in addr :", u.Agent.RemoteAddr())
}

func RemoveUser(a network.AgentInf) {
	mu.Lock()
	delete(Users, a)
	mu.Unlock()
	log.Debug("client out addr :", a.RemoteAddr())
}
