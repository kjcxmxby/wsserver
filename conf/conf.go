package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var SvrConf struct {
	WSAddr          string
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	HTTPTimeout     int
	LogLvl          int
	LogPath         string
}

func init() {

	conf_data, err := ioutil.ReadFile("./conf/conf.json")

	if err != nil {
		log.Fatalf("init conf read file error:%v", err)
	}

	err = json.Unmarshal(conf_data, &SvrConf)

	if err != nil {
		log.Fatalf("init conf unmarshal error : %v", err)
	}
}
