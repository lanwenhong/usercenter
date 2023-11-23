package user

import (
	"context"
	"fmt"
	"net/http"
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
	//Email string `form:"email" binding:"omitempty,max=128,email"`
	Email string `form:"email" binding:"EmailValidator,max=128" reg_error_info:"邮箱格式错误"`
	//Email    string `form:"email" binding:"EmailValidator"`
	Mobile   string `form:"mobile" binding:"omitempty,max=11"`
	Username string `form:"username" binding:"omitempty,max=128"`
	Password string `form:"password" binding:"required,min=6,max=10"`
}

func (uch *UserHandler) Login(ctx context.Context) error {
	ul := UserLoginData{}
	if err := uch.C.ShouldBindWith(&ul, binding.Query); err != nil {
		logger.Warnf(ctx, "login data err: %s", err.Error())
		//err := respcode.NewUserCenterErr(respcode.ERR, "参数错误")
		err := respcode.NewUserCenterErr(respcode.ERR, ut.ValidatErr(ul, err))
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
	}
	if ul.Email == "" && ul.Mobile == "" && ul.Username == "" {
		err := respcode.NewUserCenterErr(respcode.ERR, "username/email/mobile must have one")
		logger.Warnf(ctx, "err: %s", err.Error())
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
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
		err := respcode.NewUserCenterErr(respcode.ERR_USER, "username or password error")
		if ret.Error != nil {
			logger.Warnf(ctx, "query db users: %s", ret.Error.Error())
		}
		logger.Debugf(ctx, "ret.RowsAffected: %d", ret.RowsAffected)
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return ret.Error
	}
	px := strings.Split(user.Password, "$")
	pass_enc, _ := ut.CreatePassWithRand(ul.Password, px[1])
	if pass_enc != user.Password {
		err := respcode.NewUserCenterErr(respcode.ERR_AUTH, "username or password error")
		logger.Warnf(ctx, "check pass err: %s", err.Error())
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
	}
	if user.Status != respcode.STATUS_OK {
		//err := errors.New("status error")
		err := respcode.NewUserCenterErr(respcode.ERR_AUTH, "status error")
		logger.Warnf(ctx, "check status err: %s", err.Error())
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
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
		err := respcode.NewUserCenterErr(respcode.ERR_DB, "")
		resp := respcode.RespError[string](err.Code, err.Msg, "", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
	}
	uch.C.SetCookie("sid", sek, 3600, "/", "127.0.0.1", false, true)
	resp := respcode.RespSucc[map[string]interface{}](respcode.OK, userinfo)
	uch.C.JSON(http.StatusOK, resp)
	return nil
}
