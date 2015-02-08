package models

type ReturnDto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
