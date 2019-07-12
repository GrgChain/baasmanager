package model

type Dashboard struct {
	Users      int64 `json:"users"`
	Chains     int64 `json:"chains"`
	Channels   int64 `json:"channels"`
	Chaincodes int64 `json:"chaincodes"`
}


type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}