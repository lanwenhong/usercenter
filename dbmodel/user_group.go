package dbmodel

type UserGroup struct {
	ID      uint64 `gorm:"pcolumn:id;primaryKey;`
	Userid  uint64 `gorm:"column:userid;type:bigint;"`
	Groupid uint64 `gorm:"column:groupid;type:bigint;"`
	Ctime   uint64 `gorm:"column:ctime;type:int unsigned;`
	Utime   uint64 `gorm:"column:utime;type:int unsigned;`
}
