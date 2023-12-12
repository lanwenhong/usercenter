package user

import (
	"context"
	"strconv"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/session"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

type GroupAdddData struct {
	GroupId int64 `form:"userid" binding:"required" reg_error_info:"userid格式错误"`
	Userid  int64 `form:"groupid" binding:"required" reg_error_info:"userid格式错误"`
}

func (uch *UserHandler) AddGroup(ctx context.Context) error {
	gad := GroupAdddData{}
	if err := uch.C.ShouldBindWith(&gad, binding.Form); err != nil {
		logger.Debugf(ctx, "addgroup binding data: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(gad, err), "", "")
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
	userid = gad.Userid

	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	createid, _ := util.Genid(ctx, db)
	t := uint64(time.Now().Unix())
	user_group := dbmodel.UserGroup{
		ID:     createid,
		Userid: uint64(userid),
		Ctime:  t,
		Utime:  t,
	}
	ret := db.WithContext(ctx).Create(&user_group)
	if ret.Error != nil {
		logger.Warnf(ctx, "user group add: %s", ret.Error.Error())
		return respcode.RetError[string](uch.C, respcode.ERR_DB, "user group add err", "", "")
	}
	q_user_group := map[string]interface{}{}
	ret = db.WithContext(ctx).Table("users").Where("id=?", createid).Find(&q_user_group)
	if ret.Error != nil {
		logger.Warnf(ctx, "get user group relation: %s", ret.Error.Error())
		return respcode.RetError[string](uch.C, respcode.ERR_DB, "get user group relation err", "", "")
	}
	ret_ug := map[string]interface{}{}
	keys := []string{
		"id",
		"userid",
		"groupid",
	}
	for _, k := range keys {
		v, _ := q_user_group[k]
		ret_ug[k] = strconv.FormatUint(v.(int64), 10)
	}
	return respcode.RetSucc[map[string]interface{}](uch.C, ret_ug)
}
