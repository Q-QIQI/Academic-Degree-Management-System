package entity

type Org struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Publickey   string `json:"publickey" xorm:"VARCHAR(100)"`
	Channels    string `json:"channels" xorm:"VARCHAR(100)"`
	Chain       int `json:"chain" xorm:"INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(100)"`
}
