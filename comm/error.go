package comm

import (
	"errors"
)

var (
	MSG_ERROR                = errors.New("MSG_ERROR")
	MSG_NO_REGISTER          = errors.New("MSG_NO_REGISTER")
	MSG_DATA_ERROR           = errors.New("MSG_DATA_ERROR")
	MSG_ALREADY_REGISTER     = errors.New("MSG_ALREADY_REGISTER")
	MSG_REGISTER_PARAM_ERROR = errors.New("MSG_REGISTER_PARAM_ERROR")
)
