package user

import (
	"context"
	"errors"
	"net/http"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/session"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

type UserRejsterData struct {
	Email    string `form:"email" binding:"required,email"`
	Mobile   string `form:"mobile" binding:"required,len=11"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required,min=6,max=10"`
}

func (uch *UserHandler) Register(ctx context.Context) error {
	uj := UserRejsterData{}
	if err := uch.C.ShouldBind(&uj); err != nil {
		logger.Warnf(ctx, "bind data err: %s", err.Error())
		resp := respcode.RespError[string](respcode.ERR, err.Error(), "参数错误", "")
		uch.C.JSON(http.StatusOK, resp)
		return err
	}
	logger.Debugf(ctx, "email: %s mobile: %s username: %s password: %s", uj.Email, uj.Mobile, uj.Username, uj.Password)
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	user := dbmodel.Users{}
	ret := db.WithContext(ctx).Select("id").Where("mobile=? or username=? or email=?", uj.Mobile, uj.Username, uj.Email).Find(&user)

	if ret.Error != nil {
		//create user
		logger.Warnf(ctx, "db query %s", ret.Error.Error())
		resp := respcode.RespError[string](respcode.ERR, ret.Error.Error(), "", "")
		uch.C.JSON(http.StatusOK, resp)
		return ret.Error
	}
	if ret.RowsAffected == 0 {
		createid, _ := util.Genid(ctx, db)
		enc_pass, _ := ut.CreatePassword(uj.Password)
		nuser := dbmodel.Users{
			Username: uj.Username,
			Mobile:   uj.Mobile,
			Email:    uj.Email,
			Password: enc_pass,
			ID:       createid,
			Ctime:    uint64(time.Now().Unix()),
			Status:   respcode.STATUS_OK,
		}
		ret := db.WithContext(ctx).Create(&nuser)
		if ret.Error != nil {
			logger.Warnf(ctx, "insert db err:%s", ret.Error.Error())
			resp := respcode.RespError[string](respcode.ERR, ret.Error.Error(), "", "")
			uch.C.JSON(http.StatusOK, resp)
			return ret.Error
		} else {
			user, _ := uch.getUser(ctx, createid)
			resp := respcode.RespSucc[map[string]interface{}](respcode.OK, user)
			if uch.Cookie == "" {
				sek := util.GenKsuid()
				se := session.NewSession(sek)
				sdata := map[string]interface{}{
					"userid":   user["id"].(int64),
					"username": user["username"].(string),
					"isadmin":  user["isadmin"].(int64),
					"status":   user["status"].(int64),
				}
				se.Set(ctx, sek, sdata, time.Duration(3600))
				uch.C.SetCookie("sid", sek, 3600, "/", "127.0.0.1", false, true)
			}
			uch.C.JSON(http.StatusOK, resp)
			return nil
		}
	}
	//Found user data
	err := errors.New("username or email or mobile exist")
	resp := respcode.RespError[string](respcode.ERR, err.Error(), "", "")
	uch.C.JSON(http.StatusOK, resp)
	return err
}
