package controller

import (
	"net/http"
	"strconv"
	"wsserver/log"
	"wsserver/msg"
	"wsserver/user"
)

func broadcast(msg interface{}) {
	for _, u := range user.Users {

		log.Debug("broadcast conn :", u.Agent.RemoteAddr())
		log.Debug("broadcast msg:", msg)

		u.Agent.WriteMsg(msg)
	}
}

func LessonStatusUpdate(req *http.Request, w http.ResponseWriter) {

	lesson_id_str := req.PostFormValue("lesson_id")
	status_str := req.PostFormValue("status")

	log.Debug("Get Req LessonStatusUpdate lesson_id:", lesson_id_str, ",status:", status_str)

	lesson_id, err := strconv.Atoi(lesson_id_str)
	if err != nil {
		log.Error("StatusController LessonStatusUpdate param lesson_id atoi error:", err)
		return
	}

	status, err := strconv.Atoi(status_str)
	if err != nil {
		log.Error("StatusController LessonStatusUpdate param status atoi error:", err)
		return
	}

	var data msg.Data
	data.Lesson_id = lesson_id
	data.Status = status

	var pushdata msg.PushData
	pushdata.Action = "/websocket/v4/ibl/lesson/list/status"
	pushdata.Data = data
	pushdata.Err = ""
	pushdata.Status = 1001

	broadcast(&pushdata)

	w.Write([]byte("{status:1001,data:1,err:"))
}
