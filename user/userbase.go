package user

import (
	"context"
	"reflect"
	"strconv"
	"time"
	"usercenter/respcode"
	"usercenter/sruntime"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

type BaseObjectHandler struct {
	Qdata map[string]interface{}
	C     *gin.Context
	Table string
}

type BaseOpFunc func(ctx, context.Context) error

func (boh *BaseObjectHandler) Insert(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		boh.Qdata["ctime"] = uint64(time.Now().Unix())
		createid, _ := util.Genid(ctx, db)
		boh.Qdata["id"] = createid
		ret := db.WithContext(ctx).Create(&boh.Qdata)
		if ret.Error != nil {
			logger.Warnf(ctx, "insert: %s", ret.Error.Error())
			return nil, respcode.ERR_DB
		}
		ret := map[string]interface{}{
			"id": createid,
		}
		return ret, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Update(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		xid, _ := boh.Qdata["id"]
		xt := reflect.ValueOf(xid)
		delete(boh.Qdata, "id")
		switch xt {
		case reflect.Slice:
			xids := xid.([]uint64)
			ret := db.WithContext(ctx).Table(boh.Table).Where("id in (?)", xids).Updates(boh.Qdata)
			if ret.Error != nil {
				logger.Debugf(ctx, "update: %s", ret.Error.Error())
				return nil, respcode.ERR_DB
			}
			boh.Qdata["_rows"] = ret.RowsAffected
		case reflect.Uint64:
			sid := xid.(uint64)
			ret := db.WithContext(ctx).Table(boh.Table).Where("id = ?", sid).Updates(boh.Qdata)
			if ret.Error != nil {
				logger.Debugf(ctx, "update: %s", ret.Error.Error())
				return nil, respcode.ERR_DB
			}
			boh.Qdata["id"] = strconv.FormatUint(sid, 10)
		default:
			return nil, respcode.ERR_PARAM

		}
		return boh.Qdata, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Delete(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		xid, _ := boh.Qdata["id"]
		id := xid.(uint64)
		ret := db.WithContext(ctx).Table(boh.Table).Where("id = ?", id).Delete(map[string]interface{})
		if ret.Error != nil {
			return nil, respcode.ERR_DB
		}
		return map[string]interface{}{}, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) convertRow(bdata map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	nk := []string{
		"id",
		"userid",
		"groupid",
		"roleid",
		"permid",
		"parenid",
	}
	for k, v := range bdata {
		if ut.Foundin[string](nk, k) {
			ret[k] = strconv.FormatUint(v, 10)
		}
		if k == "ctime" || k == "utime" {
			tt, _ := ret[k]
			tm := time.Unix(int64(tt), 0)
			ret[k] = tm.Format("2006-01-02 15:04:05")
		}
	}
	return ret
}

func (boh *BaseObjectHandler) convertRows(bdata []map[string]interface{}) []map[string]interface{} {
	ret := make([]map[string]interface{})
	for row, _ := range bdata {
		new_one := boh.convertRow(row)
		ret = append(ret, new_one)
	}
	return ret
}

func (boh *BaseObjectHandler) GetData(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		xid, _ := boh.Qdata["id"]
		id := xid.(uint64)
		ddata := map[string]interface{}{}
		db.WithContext(ctx).Table(boh.Table).Where("id = ?", id).Find(&ddata)
		if ret.Error != nil {
			return nil, respcode.ERR_DB
		}
		ret := boh.convertRwo(ddata)
		return ret, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Post(ctx context.Context) error {
	var e_code int = respcode.OK
	data := map[string]interface{}{}
	q_str := c.Param("basequery")
	switch q_str {
	case "add":
		data, e_code = boh.Insert(ctx)
	case "mod":
		data, e_code = boh.Update(ctx)

	}
	if e_code != respcode.OK {
		return respcode.RetError[string](boh.C, e_code, "", "", "")
	}
	return respcode.RetSucc[map[string]interface{}](boh.C, data)
}

func (boh *BaseObjectHandler) Get(ctx context.Context) error {
	return nil
}
