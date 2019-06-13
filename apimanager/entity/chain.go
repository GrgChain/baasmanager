package entity

type Chain struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"not null VARCHAR(64)"`
	UserAccount string `json:"userAccount" xorm:"not null VARCHAR(100)"`
	Description string `json:"description" xorm:"VARCHAR(255)"`
	Consensus   string `json:"consensus" xorm:"not null VARCHAR(10)"`
	PeersOrgs   string `json:"peersOrgs" xorm:"not null VARCHAR(100)"`
	OrderCount  int    `json:"orderCount" xorm:"not null INT(11)"`
	PeerCount   int    `json:"peerCount" xorm:"not null INT(11)"`
	TlsEnabled  string `json:"tlsEnabled" xorm:"not null VARCHAR(5)"`
	Status      int    `json:"status" xorm:"default 0 INT(11)"` //0定义 1已构建 2运行中 3已停止
	Created     int64  `json:"created" xorm:"not null BIGINT(20)"`
}
