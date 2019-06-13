package model

type RespData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

type Dashboard struct {
	Users      int64 `json:"users"`
	Chains     int64 `json:"chains"`
	Channels   int64 `json:"channels"`
	Chaincodes int64 `json:"chaincodes"`
}
