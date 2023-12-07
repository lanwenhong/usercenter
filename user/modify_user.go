package user

import (
	"context"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/session"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type ModifyUserData struct {
	Mobile   string `form:"mobile" binding:"omitempty,max=11"  reg_error_info:"手机格式错误"`
	Username string `form:"username" binding:"omitempty,max=128" reg_error_info:"用户名格式错误"`
	Password string `form:"password" binding:"omitempty,min=6,max=10" reg_error_info:"密码格式错误"`
	Status   int    `form:"status" binding:"omitempty,len=1" reg_error_info:"用户状态格式错误"`
	Extend   string `form:"extend" binding:"omitempty,max=8192" reg_error_info:"扩展字段格式错误"`
	Userid   int64  `form:"userid" binding:"required" reg_error_info:"userid格式错误"`
}

func (uch *UserHandler) ModifyUser(ctx context.Context) error {
	mud := ModifyUserData{}
	if err := uch.C.ShouldBindWith(&mud, binding.Form); err != nil {
		//if err := uch.C.ShouldBind(&mud); err != nil {
		logger.Debugf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(mud, err), "", "")
	}
	se := session.NewSession(uch.Cookie)
	it_isadmin, _ := se.GetData(ctx, "isadmin")
	it_userid, _ := se.GetData(ctx, "userid")

	isadmin := int64(it_isadmin.(float64))
	userid := int64(it_userid.(float64))
	logger.Debugf(ctx, "userid=%d,isadmin=%d", userid, isadmin)
	if isadmin == 0 {
		logger.Warnf(ctx, "permission deny")
		return respcode.RetError[string](uch.C, respcode.ERR_PARAM, "permission deny", "", "")
	}
	userid = mud.Userid

	up_user := map[string]interface{}{}
	keys := []string{
		"username",
		"password",
		"extend",
		"status",
		"mobile",
	}
	for _, k := range keys {
		switch k {
		case "username":
		case "extend":
		case "status":
		case "mobile":
			{
				logger.Debugf(ctx, "k=%s", k)
				v := uch.C.PostForm(k)
				logger.Debugf(ctx, "v=%s", v)
				if v != "" {
					up_user[k] = v
				}
			}
		case "password":
			{
				v := uch.C.Query(k)
				if v != "" {
					enc_pass, _ := ut.CreatePassword(v)
					up_user[k] = enc_pass
				}

			}

		}
	}
	if len(up_user) == 0 {
		logger.Warnf(ctx, "no modify info")
		return respcode.RetError[string](uch.C, respcode.ERR_PARAM, "modify user error", "", "")
	}
	up_user["utime"] = uint64(time.Now().Unix())
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.WithContext(ctx).Model(&dbmodel.Users{}).Where("id = ?", userid).Updates(up_user)
	if ret.Error != nil {
		logger.Warnf(ctx, "modify user: %s", ret.Error.Error())
		return respcode.RetError[string](uch.C, respcode.ERR_DB, "modify user error", "", "")
	}
	user := map[string]interface{}{}
	ret = db.WithContext(ctx).Table("users").Select("id", "username", "mobile", "status", "ctime", "utime").Where("id=?", userid).Find(&user)
	if ret.Error != nil {
		logger.Warnf(ctx, "get user info error")
		return respcode.RetError[string](uch.C, respcode.ERR_DB, "get user info error", "", "")
	}
	return respcode.RetSucc[map[string]interface{}](uch.C, user)
}
