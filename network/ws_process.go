package network

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"wsserver/comm"
	"wsserver/log"
)

func NewWsProcess() *WsProcess {
	p := &WsProcess{msgInfo: make(map[uint32]*ProcessInfo),
		msgIds: make(map[reflect.Type]uint32),
	}

	return p
}

type Message struct {
	MsgId uint32 `json:"msg_id"`
	Data  string `json:"data"`
}

type MsgHandler func([]interface{})

type ProcessInfo struct {
	msgType reflect.Type
	handler MsgHandler
}

type WsProcess struct {
	msgInfo map[uint32]*ProcessInfo
	msgIds  map[reflect.Type]uint32
}

func (p *WsProcess) Register(msg_id uint32, msg interface{}, handler MsgHandler) error {

	if _, ok := p.msgInfo[msg_id]; ok {
		log.Error("WsProcess Register msg_id is already registered:", msg_id)
		return comm.MSG_ALREADY_REGISTER
	}

	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Error("WsProcess Register json message pointer required")
		return comm.MSG_REGISTER_PARAM_ERROR
	}

	if _, ok := p.msgIds[msgType]; ok {
		log.Error("WsProcess Register msg is already registered:", msg_id)
		return comm.MSG_ALREADY_REGISTER
	}

	i := new(ProcessInfo)
	i.msgType = msgType
	i.handler = handler

	p.msgIds[msgType] = msg_id
	p.msgInfo[msg_id] = i

	return nil
}

// must goroutine safe
func (p *WsProcess) Route(msg_id uint32, msg interface{}, userData interface{}) error {
	info, ok := p.msgInfo[msg_id]

	if !ok {
		log.Error("WsProcess Router no msg id :", msg_id)
		return comm.MSG_NO_REGISTER
	}

	if info.handler != nil {

		defer func() {
			if err := recover(); err != nil {
				log.Error("WsProcess Router Panic err:", err)
			}
		}()

		log.Debug("WsProcess Route handler begin msg:", msg)
		info.handler([]interface{}{msg, userData})
		log.Debug("WsProcess Route handler end msg:", msg)
	} else {
		log.Error("WsProcess Router no handler req msg id:", msg_id)
		return comm.MSG_NO_REGISTER
	}

	return nil
}

// must goroutine safe
func (p *WsProcess) Unmarshal(data []byte) (uint32, interface{}, error) {

	var msg Message

	err := json.Unmarshal(data, &msg)

	if err != nil {
		log.Error("WsProcess json Unmarshal error:", err)
		return 0, nil, comm.MSG_ERROR
	}

	info, ok := p.msgInfo[msg.MsgId]

	if !ok {
		log.Error("WsProcess error msg id:", msg.MsgId)
		return 0, nil, comm.MSG_NO_REGISTER
	}

	b, err := base64.StdEncoding.DecodeString(msg.Data)

	if err != nil {
		log.Error("WsProcess base64 decode error:", err)
		return 0, nil, comm.MSG_DATA_ERROR
	}

	log.Debug("WsProcess base64 decode data:", string(b))

	data_msg := reflect.New(info.msgType.Elem()).Interface()

	err = json.Unmarshal(b, data_msg)

	if err != nil {
		log.Error("WsProcess msg data Unmarshal error:", err)
		return 0, nil, comm.MSG_DATA_ERROR
	}

	return msg.MsgId, data_msg, nil
}

// must goroutine safe
func (p *WsProcess) Marshal(msg interface{}) ([]byte, error) {

	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Error("WsProcess Marshal json message pointer required")
		return nil, comm.MSG_ERROR
	}

	msg_id, ok := p.msgIds[msgType]

	if !ok {
		log.Error("WsProcess Marshal msg no register:", msgType)
		return nil, comm.MSG_NO_REGISTER
	}

	d, err := json.Marshal(msg)

	if err != nil {
		log.Error("WsProcess Marshal msg json marshal error: ", err)
		return nil, comm.MSG_ERROR
	}

	//b_str := base64.StdEncoding.EncodeToString(d)

	// data
	//m := map[string]interface{}{"msg_id": msg_id, "data": d}

	var rsp_msg Message

	rsp_msg.MsgId = msg_id
	rsp_msg.Data = string(d)

	data, err := json.Marshal(rsp_msg)

	return data, err
}
