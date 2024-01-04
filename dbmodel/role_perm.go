package dbmodel

type RolePerm struct {
	ID     uint64  `gorm:"pcolumn:id;primaryKey;`
	Roleid uint64  `gorm:"column:roleid;type:bigint;"`
	Permid uint64  `gorm:"column:permid;type:bigint;"`
	Ctime  *uint64 `gorm:"column:ctime;type:int unsigned;`
	Utime  *uint64 `gorm:"column:utime;type:int unsigned;`
}
