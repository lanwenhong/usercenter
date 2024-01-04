package user

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"gorm.io/gorm"
)

type RolesAddData struct {
	Name string `form:"name" binding:"required" reg_error_info:"角色名字格式错误"`
	Info string `form:"info" binding:"omitempty" reg_error_info:"角色信息格式错误"`
}

type RolesDelData struct {
	Id uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
}

type RolesDelPermData struct {
	Id         uint64   `form:"id" binding:"omitempty" reg_error_info:"id格式错误"`
	Permid     []uint64 `form:"permid" binding:"omitempty" reg_error_info:"permid格式错误"`
	RolePermid []uint64 `form:"roleperm_id" binding:"omitempty" reg_error_info:"roleperm_id格式错误"`
}

type RolesModDataWithId struct {
	Id   uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
	Info string `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name string `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
}

type RolesModDataWithIds struct {
	Id   []uint64 `form:"id" binding:"required" reg_error_info:"ids格式错误"`
	Info string   `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name string   `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
}

type RoleData struct {
	Id uint64 `form:"id" binding:"omitempty" reg_error_info:"id格式错误"`
}

type RolesAddPermsData struct {
	Id     uint64   `form:"id" binding:"required" reg_error_info:"id格式错误"`
	Permid []uint64 `form:"permid" binding:"required" reg_error_info:"permid格式错误"`
}

type RolePermDataList struct {
	Id       []uint64  `form:"id" binding:"omitempty" reg_error_info:"ids格式错误"`
	Info     []string  `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name     []string  `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
	Fctime   time.Time `form:"fctime" binding:"omitempty" time_format:"2006-01-02 15:04:05" reg_error_info:"ctime起始格>式错误"`
	Tctime   time.Time `form:"tctime" binding:"omitempty" time_format:"2006-01-02 15:04:05" reg_error_info:"ctime结束格>式错误"`
	Page     int       `form:"page" binding:"required" reg_error_info:"page格式错误"`
	PageSize int       `form:"page_size" binding:"required" reg_error_info:"page_size格式错误"`
}

type RolesOpHandler struct {
	BaseObjectHandler
	BaseOpFuncIndex map[string]BaseOpFunc
}

func (roh *RolesOpHandler) GetDataList(ctx context.Context) (map[string]interface{}, int) {
	retdata, code := roh.BaseObjectHandler.GetDataList(ctx)
	if code != respcode.OK {
		return retdata, code
	}
	data, _ := retdata["data"]
	ldata := data.([]map[string]interface{})
	roleids := []int64{}
	for _, row := range ldata {
		id, _ := row["id"]
		sid := id.(string)
		cid, _ := strconv.ParseInt(sid, 10, 64)
		roleids = append(roleids, cid)
	}
	logger.Debugf(ctx, "roleids: %v", roleids)

	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]

	retdic := []map[string]interface{}{}
	sql := "select p.id as id,p.name as name,p.info as info,rp.roleid as roleid from perms p, role_perm rp where rp.roleid in ? and rp.permid=p.id"
	ret := db.WithContext(ctx).Raw(sql, roleids).Scan(&retdic)
	if ret.Error != nil {
		logger.Warnf(ctx, "insert: %s", ret.Error.Error())
		return nil, respcode.ERR_DB
	}
	logger.Debugf(ctx, "retdic: %v", retdic)
	if len(retdic) > 0 {
		roledict := map[string][]map[string]interface{}{}
		for _, row := range retdic {
			id, _ := row["roleid"]
			s_id := fmt.Sprintf("%d", id)
			delete(row, "roleid")
			perm, ok := roledict[s_id]
			if !ok {
				plist := []map[string]interface{}{}
				plist = append(plist, row)
				roledict[s_id] = plist
			} else {
				perm = append(perm, row)
				roledict[s_id] = perm
			}
		}
		for _, old_row := range ldata {
			id, _ := old_row["id"]
			sid := id.(string)
			if perm, ok := roledict[sid]; ok {
				old_row["perm"] = perm
			} else {
				old_row["perm"] = []map[string]interface{}{}
			}
		}
		logger.Debugf(ctx, "roledict: %v", roledict)
	} else {
		perm := []map[string]interface{}{}
		for _, row := range ldata {
			row["perm"] = perm
		}
	}
	return retdata, respcode.OK
}

func (roh *RolesOpHandler) QlistOpFunc(ctx context.Context) error {
	rpdl := RolePermDataList{}
	if err := roh.C.ShouldBindWith(&rpdl, binding.Query); err != nil {
		logger.Warnf(ctx, "qlist binding data: %s", err.Error())
		return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rpdl, err), "", "")
	}
	logger.Debugf(ctx, "qlist handler")
	roh.Qdata, _ = ut.Stru2Map(ctx, rpdl)
	return roh.Get(ctx)
}

func (roh *RolesOpHandler) DelRolePermFunc(ctx context.Context) error {
	rdpd := RolesDelPermData{}
	if err := roh.C.ShouldBindWith(&rdpd, binding.Form); err != nil {
		logger.Warnf(ctx, "addperm binding data: %s", err.Error())
		return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rdpd, err), "", "")
	}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	var ret *gorm.DB = nil
	if rdpd.Id > 0 && len(rdpd.Permid) > 0 {
		where := "roleid = ? and "
		where2 := ""
		if len(rdpd.Permid) == 1 {
			where2 = "permid = ?"
			where = where + where2
			ret = db.WithContext(ctx).Table("role_perm").Where(where, rdpd.Id, rdpd.Permid[0]).Delete(map[string]interface{}{})
		} else {
			where2 = "permid in ?"
			where = where + where2
			ret = db.WithContext(ctx).Table("role_perm").Where(where, rdpd.Id, rdpd.Permid).Delete(map[string]interface{}{})
		}
	} else if len(rdpd.RolePermid) > 0 {
		if len(rdpd.RolePermid) == 1 {
			ret = db.WithContext(ctx).Table("role_perm").Where("id = ?", rdpd.RolePermid[0]).Delete(map[string]interface{}{})
		} else {
			ret = db.WithContext(ctx).Table("role_perm").Where("id in ?", rdpd.RolePermid).Delete(map[string]interface{}{})
		}
	}
	if ret.Error != nil {
		logger.Warnf(ctx, "insert: %s", ret.Error.Error())
		return respcode.RetError[string](roh.C, respcode.ERR_DB, "", "", "")
	}
	data := map[string]interface{}{
		"rows": ret.RowsAffected,
	}
	return respcode.RetSucc[map[string]interface{}](roh.C, data)
}

func (roh *RolesOpHandler) AddPermsFunc(ctx context.Context) error {
	rapd := RolesAddPermsData{}
	if err := roh.C.ShouldBindWith(&rapd, binding.Form); err != nil {
		logger.Warnf(ctx, "addperm binding data: %s", err.Error())
		return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rapd, err), "", "")
	}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	rp := dbmodel.RolePerm{}
	for _, pid := range rapd.Permid {
		createid, _ := util.Genid(ctx, db)
		rp.ID = createid
		rp.Roleid = rapd.Id
		rp.Permid = pid
		ctime := uint64(time.Now().Unix())
		utime := uint64(time.Now().Unix())
		rp.Ctime = &ctime
		rp.Utime = &utime
		ret := db.WithContext(ctx).Table("role_perm").Create(&rp)
		if ret.Error != nil {
			logger.Warnf(ctx, "insert: %s", ret.Error.Error())
			return respcode.RetError[string](roh.C, respcode.ERR_DB, "", "", "")
		}
	}
	rows := len(rapd.Permid)
	data := map[string]interface{}{
		"rows": rows,
	}
	return respcode.RetSucc[map[string]interface{}](roh.C, data)
}

func (roh *RolesOpHandler) AddOpFunc(ctx context.Context) error {
	rad := RolesAddData{}
	if err := roh.C.ShouldBindWith(&rad, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rad, err), "", "")
	}
	roh.Qdata, _ = ut.Stru2Map(ctx, rad)
	return roh.Post(ctx)
}

func (roh *RolesOpHandler) DelOpFunc(ctx context.Context) error {
	rdd := RolesDelData{}
	if err := roh.C.ShouldBindWith(&rdd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rdd, err), "", "")
	}
	roh.Qdata, _ = ut.Stru2Map(ctx, rdd)
	logger.Debugf(ctx, "%v", roh.Qdata)
	return roh.Post(ctx)
}

func (roh *RolesOpHandler) ModOpFunc(ctx context.Context) error {
	pids := roh.C.PostFormArray("id")
	logger.Debugf(ctx, "pids: %v", pids)
	if len(pids) == 1 {
		rmd := RolesModDataWithId{}
		if err := roh.C.ShouldBindWith(&rmd, binding.Form); err != nil {
			logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
			return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rmd, err), "", "")
		}
		roh.Qdata, _ = ut.Stru2Map(ctx, rmd)
		return roh.Post(ctx)
	} else if len(pids) > 1 {
		rmd := RolesModDataWithIds{}
		if err := roh.C.ShouldBindWith(&rmd, binding.Form); err != nil {
			logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
			return respcode.RetError[string](roh.C, respcode.ERR, ut.ValidatErr(rmd, err), "", "")
		}
		roh.Qdata, _ = ut.Stru2Map(ctx, rmd)
		return roh.Post(ctx)
	}
	return respcode.RetError[string](roh.C, respcode.ERR, "modify data err", "", "")
}

func RolesOpHandlerNew(c *gin.Context, cookie string) *RolesOpHandler {
	roh := RolesOpHandler{}
	roh.C = c
	roh.Cookie = cookie
	roh.Table = "roles"

	roh.BaseOpFuncIndex = map[string]BaseOpFunc{
		"add":    roh.AddOpFunc,
		"mod":    roh.ModOpFunc,
		"delete": roh.DelOpFunc,
		"list":   roh.QlistOpFunc,
		//"q":       roh.QopFunc,
		"addperm": roh.AddPermsFunc,
		"delperm": roh.DelRolePermFunc,
	}
	roh.InsertFunc = roh.Insert
	roh.UpdateFunc = roh.Update
	roh.DeleteFunc = roh.Delete
	roh.QlistFunc = roh.GetDataList
	roh.Qfunc = roh.GetData
	return &roh
}

func RoleOp(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_edit")
	roh := RolesOpHandlerNew(c, cookie.(string))
	logger.Debugf(ctx, "grouped: %s", grouped)
	if op, ok := roh.BaseOpFuncIndex[grouped]; ok {
		op(ctx)
	} else {
		respcode.RetError[string](c, respcode.ERR, "not found method", "", "")
	}
}

func RoleQuery(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_query")
	logger.Debugf(ctx, "grouped: %s", grouped)
	roh := RolesOpHandlerNew(c, cookie.(string))
	if op, ok := roh.BaseOpFuncIndex[grouped]; ok {
		op(ctx)
	} else {
		respcode.RetError[string](c, respcode.ERR, "not found method", "", "")
	}
}
