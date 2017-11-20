package msg

type LoginReq struct {
	Sid string `json:"sid"`
}

type LoginRsp struct {
	Status uint32 `json:"status"`
	Err    string `json:"err"`
}
