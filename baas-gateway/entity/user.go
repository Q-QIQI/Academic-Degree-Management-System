package entity

type User struct {
	Id       int    `json:"id"       xorm:"not null pk autoincr INT(11)"`
	Account  string `json:"account"  xorm:"not null unique VARCHAR(30)"`
	Password string `json:"password" xorm:"not null VARCHAR(100)"`
	Phone    string `json:"phone"    xorm:"VARCHAR(200)"`
	Name     string `json:"name"     xorm:"not null VARCHAR(200)"`
	IdCard   string `json:"idCard"   xorm:"VARCHAR(200)"`
	Org      string `json:"org"      xorm:"-"`
	Type     int    `json:"type"     xorm:"not null INT(11)"`
	Created  int64  `json:"created"  xorm:"not null BIGINT(20)"`
	Updated  string `json:"updated"  xorm:"VARCHAR(200)"`
	IsOk     string `json:"isok"  xorm:"VARCHAR(200)"`
	Img      string `json:"img"  xorm:"VARCHAR(200)"`

	// 用户公钥
	// PubKey string `json:"pubKey" xorm:"VARCHAR(600)"`
}

type UserInfo struct {
	Id           int      `json:"id"`
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
	Account      string   `json:"account"`
	Type         int      `json:"type"`
}

type UserRole struct {
	UserId  int    `json:"userId"  xorm:"not null pk INT(11)"`
	RoleKey string `json:"roleKey" xorm:"not null pk VARCHAR(20)"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserDetail struct {
	User
	Roles []string `json:"roles"`
}
