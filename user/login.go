package user

import (
	"context"
	"fmt"
	"strings"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/session"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type UserLoginData struct {
	Email string `form:"email" binding:"EmailValidator,max=128" reg_error_info:"邮箱格式错误"`
	Mobile   string `form:"mobile" binding:"omitempty,len=11" reg_error_info:"手机格式错误"`
	Username string `form:"username" binding:"omitempty,max=128" reg_error_info:"用户名格式错误"`
	Password string `form:"password" binding:"required,min=6,max=10" reg_error_info:"密码格式错误"`
}

func (uch *UserHandler) Login(ctx context.Context) error {
	ul := UserLoginData{}
	if err := uch.C.ShouldBindWith(&ul, binding.Query); err != nil {
		logger.Warnf(ctx, "login data err: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(ul, err), "", "")
	}
	if ul.Email == "" && ul.Mobile == "" && ul.Username == "" {
		logger.Warnf(ctx, "err: %s", "username/email/mobile must have one")
		return respcode.RetError[string](uch.C, respcode.ERR, "username/email/mobile must have one", "", "")
	}
	log_key := ""
	log_val := ""
	if ul.Email != "" {
		log_key = "email"
		log_val = ul.Email
	}
	if ul.Mobile != "" {
		log_key = "mobile"
		log_val = ul.Mobile
	}
	if ul.Username != "" {
		log_key = "username"
		log_val = ul.Username
	}
	where_cond := fmt.Sprintf("%s = ?", log_key)
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	user := dbmodel.Users{}
	ret := db.WithContext(ctx).Select("id,username,email,password,isadmin,status").Where(where_cond, log_val).Find(&user)
	if ret.Error != nil || ret.RowsAffected == 0 {
		if ret.Error != nil {
			logger.Warnf(ctx, "query db users: %s", ret.Error.Error())
		}
		return respcode.RetError[string](uch.C, respcode.ERR_USER, "username or password error", "", "")
	}
	px := strings.Split(user.Password, "$")
	pass_enc, _ := ut.CreatePassWithRand(ul.Password, px[1])
	if pass_enc != user.Password {
		logger.Warnf(ctx, "username or password error")
		return respcode.RetError[string](uch.C, respcode.ERR_AUTH, "username or password error", "", "")
	}
	if user.Status != respcode.STATUS_OK {
		logger.Warnf(ctx, "status error")
		return respcode.RetError[string](uch.C, respcode.ERR_AUTH, "status error", "", "")

	}
	db.WithContext(ctx).Model(&user).Update("logtime", uint64(time.Now().Unix()))
	//set cookie
	userinfo, _ := uch.getUser(ctx, user.ID)
	logger.Debugf(ctx, "userinfo: %v", userinfo)
	allperms := []string{}
	_, ok := userinfo["allperm"]
	if ok {
		for _, perm := range userinfo["allperm"].([]map[string]interface{}) {
			oneperm := perm["name"].(string)
			allperms = append(allperms, oneperm)
		}
	}
	sdata := map[string]interface{}{
		"userid":   userinfo["id"].(int64),
		"username": userinfo["username"].(string),
		"isadmin":  userinfo["isadmin"].(int64),
		"status":   userinfo["status"].(int64),
		"allperm":  allperms,
	}
	se := session.NewSession(uch.Cookie)
	sek, err := se.Update(ctx, uch.Cookie, sdata, time.Duration(3600))
	if err != nil {
		logger.Warnf(ctx, "update se err: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR_DB, "update db error", "", "")
	}
	uch.C.SetCookie("sid", sek, 3600, "/", "127.0.0.1", false, true)
	return respcode.RetSucc[map[string]interface{}](uch.C, userinfo)
}
