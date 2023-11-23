package dbmodel

type VerifyCode struct {
	Id       uint64 `gorm:"primaryKey;"`
	VerifyId string `gorm:"column:verify_id;index:verify_id"`
	Answer   string `gorm:"column:answer"`
	Valid    uint   `gorm:"column:valid"`
	Stime    uint64 `gorm:"column:stime"`
	Etime    uint64 `gorm:"column:etime"`
}

func (VerifyCode) TableName() string {
	return "verify_code"
}
