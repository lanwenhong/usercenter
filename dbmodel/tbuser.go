package dbmodel

type Users struct {
	ID       uint64 `gorm:"primaryKey;`
	Username string `gorm:"column:username;index:username,unique"`
	Password string `gorm:"column:password"`
	UserType uint8  `gorm:"column:usertype;default:1"`
	Email    string `gorm:"column:email;index:email,unique"`
	Mobile   string `gorm:"column:mobile;index:mobile,unique;type:varchar(18);"`
	Head     string `gorm:"column:head;type:varchar(128);"`
	Score    uint32 `gorm:"column:score;type:int;"`
	Ctime    uint64 `gorm:"column:ctime;type:int;`
	Utime    uint64 `gorm:"column:utime;type:int;`
	Logtime  uint64 `gorm:"column:logtime;type:int;`
	Regip    string `gorm:"column:regip;type:varchar(128);"`
	Status   uint8  `gorm:"column:status;type:tinyint;default:0;"`
	Isadmin  uint8  `gorm:"column:isadmin;type:tinyint;default:0;"`
	Position string `gorm:"column:position;type:varchar(128);default:"";"`
	Extend   string `gorm:"column:extend;type:varchar(8192);"`
}
