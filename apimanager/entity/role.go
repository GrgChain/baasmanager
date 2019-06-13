package entity

type Role struct {
	Rkey        string `json:"rkey" xorm:"not null pk VARCHAR(20)"`
	Name        string `json:"name" xorm:"not null VARCHAR(40)"`
	Description string `json:"description" xorm:"VARCHAR(200)"`
}
