package network

type Processor interface {
	// must goroutine safe
	Route(msg_id uint32, msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (uint32, interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}

var Default_Processor = NewWsProcess()

func Register(msg_id uint32, msg interface{}, handler MsgHandler) error {
	return Default_Processor.Register(msg_id, msg, handler)
}
