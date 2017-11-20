package msg

type Data struct {
	Lesson_id int `json:"lesson_id"`
	Status    int `json:"status"`
}

type PushData struct {
	Action string `json:"action"`
	Data   Data   `json:"data"`
	Status int    `json:"status"`
	Err    string `json:"err"`
}
