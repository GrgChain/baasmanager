package entity

type User struct {
	Id       int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Account  string `json:"account" xorm:"not null unique VARCHAR(30)"`
	Password string `json:"password" xorm:"not null VARCHAR(100)"`
	Avatar   string `json:"avatar" xorm:"VARCHAR(200)"`
	Name     string `json:"name" xorm:"not null VARCHAR(20)"`
	Created  int64  `json:"created" xorm:"not null BIGINT(20)"`
}

type UserInfo struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
	Account      string   `json:"account"`
}

type UserRole struct {
	UserId  int    `json:"userId" xorm:"not null pk INT(11)"`
	RoleKey string `json:"roleKey" xorm:"not null pk VARCHAR(20)"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserDetail struct {
	User
	Roles []string `json:"roles"`
}
