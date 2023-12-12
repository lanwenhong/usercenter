package user

import (
	"context"
	"net/http"
	"usercenter/respcode"
	"usercenter/sruntime"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

type UserHandler struct {
	C      *gin.Context
	Cookie string
}

func (uch *UserHandler) find(ids []int64, id int64) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func (uch *UserHandler) getUserBase(ctx context.Context, id uint64) (map[string]interface{}, error) {
	user := map[string]interface{}{}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.WithContext(ctx).Table("users").Where("id=?", id).Find(&user)
	if ret.Error != nil {
		logger.Warnf(ctx, "query db err: %s", ret.Error.Error())
		return user, ret.Error
	}
	delete(user, "password")
	return user, nil
}

func (uch *UserHandler) getUser(ctx context.Context, id uint64) (map[string]interface{}, error) {
	user := []map[string]interface{}{}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.WithContext(ctx).Raw(`select id,username,email,mobile,head,score,stage, 
		FROM_UNIXTIME(ctime, '%Y-%m-%d %H:%i:%s') as ctime,FROM_UNIXTIME(utime, '%Y-%m-%d %H:%i:%s') as utime,
		FROM_UNIXTIME(logtime, '%Y-%m-%d %H:%i:%s') as logtime,
		regip,status,isadmin,extend from users where id=?`, id).Scan(&user)
	if ret.Error != nil {
		return map[string]interface{}{}, ret.Error
	}
	retuser := user[0]
	groups := []map[string]interface{}{}
	ret = db.WithContext(ctx).Raw("select g.id as id,g.name as name from user_group ug, `groups` g where g.id=ug.groupid and ug.userid=?", id).Scan(&groups)
	if ret.Error != nil {
		return retuser, ret.Error
	}
	retuser["group"] = groups
	userperm := []map[string]interface{}{}
	ret = db.WithContext(ctx).Raw("select permid,roleid from user_perm where userid=?", id).Scan(&userperm)
	if ret.Error != nil {
		return retuser, ret.Error
	}
	if ret.RowsAffected > 0 {
		roles := []int64{}
		perms := []int64{}
		logger.Debugf(ctx, "%v", userperm)
		for _, v := range userperm {
			logger.Debugf(ctx, "%v", v)
			roles = append(roles, v["roleid"].(int64))
			perms = append(perms, v["permid"].(int64))
		}
		qrole := []map[string]interface{}{}
		ret = db.WithContext(ctx).Table("roles").Select("id", "name").Where("id in (?)", roles).Scan(&qrole)
		if ret.Error != nil {
			return retuser, ret.Error
		}
		if ret.RowsAffected > 0 {
			retuser["role"] = qrole
		}
		qperms := []map[string]interface{}{}
		ret = db.WithContext(ctx).Table("perms").Select("id", "name").Where("id in (?)", perms).Scan(&qperms)
		if ret.Error != nil {
			return retuser, ret.Error
		}
		if ret.RowsAffected > 0 {
			retuser["perm"] = qperms
			retuser["allperm"] = qperms
		}
		if len(roles) > 0 {
			retperms := []map[string]interface{}{}
			ret = db.WithContext(ctx).Raw("select  p.id as id,p.name as name from perms p, role_perm rp where rp.roleid in (?) and rp.permid=p.id", roles).Scan(&retperms)
			if ret.Error != nil {
				return retuser, ret.Error
			}
			allperm_id := []int64{}
			for _, v := range qperms {
				allperm_id = append(allperm_id, v["id"].(int64))
			}
			for _, row := range retperms {
				//if row["id"]
				find := uch.find(allperm_id, row["id"].(int64))
				if !find {
					rows := retuser["allperm"].([]map[string]interface{})
					rows = append(rows, row)
				}
			}
		}
	}
	return retuser, nil
}

func (uch *UserHandler) TestgetUser(ctx context.Context, id uint64) (map[string]interface{}, error) {
	return uch.getUser(ctx, id)
}

func UserOp(c *gin.Context) {

	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")

	useredit := c.Param("useredit")
	se_check, _ := c.Get("check_session")
	if se_check.(string) == "fail" {
		resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}

	logger.Debugf(ctx, "op: %s", useredit)
	uh := UserHandler{
		C:      c,
		Cookie: cookie.(string),
	}
	switch useredit {
	case "signup":
		logger.Debugf(ctx, "register user")
		uh.Register(ctx)
	case "mod":
		logger.Debugf(ctx, "modify user")
		uh.ModifyUser(ctx)
	}
}

func UserQuery(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	userquery := c.Param("userquery")
	se_check, _ := c.Get("check_session")
	if se_check.(string) == "fail" {
		resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	logger.Debugf(ctx, "op: %s", userquery)
	uh := UserHandler{
		C:      c,
		Cookie: cookie.(string),
	}
	switch userquery {
	case "login":
		logger.Debugf(ctx, "login")
		uh.Login(ctx)
	case "get_user":
		logger.Debugf(ctx, "get_user")
		uh.GetUser(ctx)
	case "query_ids":
		logger.Debugf(ctx, "query_ids")
		uh.QueryIds(ctx)

		/*case "mod":
		logger.Debugf(ctx, "modify user")
		uh.ModifyUser(ctx)*/

	}

}
