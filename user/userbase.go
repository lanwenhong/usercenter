package user

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"usercenter/respcode"
	"usercenter/sruntime"
	ut "usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

type BaseFunc func(context.Context) (map[string]interface{}, int)
type BaseObjectHandler struct {
	Qdata      map[string]interface{}
	C          *gin.Context
	Table      string
	Cookie     string
	InsertFunc BaseFunc
	UpdateFunc BaseFunc
	DeleteFunc BaseFunc
	QlistFunc  BaseFunc
	Qfunc      BaseFunc
}

type BaseOpFunc func(context.Context) error

func (boh *BaseObjectHandler) Insert(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		boh.Qdata["ctime"] = uint64(time.Now().Unix())
		createid, _ := util.Genid(ctx, db)
		boh.Qdata["id"] = createid
		logger.Debugf(ctx, "insert data: %v", boh.Qdata)
		ret := db.WithContext(ctx).Table(boh.Table).Create(&boh.Qdata)
		if ret.Error != nil {
			logger.Warnf(ctx, "insert: %s", ret.Error.Error())
			return nil, respcode.ERR_DB
		}
		re := map[string]interface{}{
			"id": createid,
		}
		return re, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Update(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		indata := map[string]interface{}{}
		for k, v := range boh.Qdata {
			if k != "id" && k != "ctime" {
				indata[k] = v
			}
		}
		indata["utime"] = uint64(time.Now().Unix())

		id, _ := boh.Qdata["id"]
		switch id.(type) {
		case uint64:
			gid := id.(uint64)
			ret := db.WithContext(ctx).Table(boh.Table).Where("id = ?", gid).Updates(indata)
			if ret.Error != nil {
				logger.Debugf(ctx, "update: %s", ret.Error.Error())
				return nil, respcode.ERR_DB
			}
			indata["id"] = strconv.FormatUint(gid, 10)
			indata["_rows"] = ret.RowsAffected
			return indata, respcode.OK
		case []uint64:
			gids := id.([]uint64)
			logger.Debugf(ctx, "%v", indata)
			ret := db.WithContext(ctx).Table(boh.Table).Where("id in (?)", gids).Updates(indata)
			if ret.Error != nil {
				logger.Debugf(ctx, "update: %s", ret.Error.Error())
				return nil, respcode.ERR_DB
			}
			ids := []string{}
			for _, id := range gids {
				sid := strconv.FormatUint(id, 10)
				ids = append(ids, sid)
			}
			indata["_rows"] = ret.RowsAffected
			indata["ids"] = ids
			return indata, respcode.OK

		}
		return nil, respcode.ERR_PARAM
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Delete(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		xid, _ := boh.Qdata["id"]
		id := xid.(uint64)
		ret := db.WithContext(ctx).Table(boh.Table).Where("id = ?", id).Delete(map[string]interface{}{})
		if ret.Error != nil {
			return nil, respcode.ERR_DB
		}
		return map[string]interface{}{}, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) convertRow(ctx context.Context, bdata map[string]interface{}) map[string]interface{} {
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
			ret[k] = strconv.FormatUint(uint64(v.(int64)), 10)
		} else {
			logger.Debugf(ctx, "k=%s v=%v", k, v)
			if k == "ctime" || k == "utime" {
				tt, _ := bdata[k]
				if tt != nil {
					logger.Debugf(ctx, "%v", tt)
					tm := time.Unix(tt.(int64), 0)
					ret[k] = tm.Format("2006-01-02 15:04:05")
				}
			} else {
				ret[k] = v
			}
		}
	}
	logger.Debugf(ctx, "ret: %v", ret)
	return ret
}

func (boh *BaseObjectHandler) convertRows(ctx context.Context, bdata []map[string]interface{}) []map[string]interface{} {
	ret := []map[string]interface{}{}
	for _, row := range bdata {
		new_one := boh.convertRow(ctx, row)
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
		ret := db.WithContext(ctx).Table(boh.Table).Where("id = ?", id).Find(&ddata)
		if ret.Error != nil {
			logger.Warnf(ctx, "query data: %s", ret.Error.Error())
			return nil, respcode.ERR_DB
		}
		re := boh.convertRow(ctx, ddata)
		return re, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) GetDataList(ctx context.Context) (map[string]interface{}, int) {
	if boh.Qdata != nil {
		db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
		db = db.WithContext(ctx).Table(boh.Table)

		logger.Debugf(ctx, "%v", boh.Qdata)
		for k, v := range boh.Qdata {
			if k == "page" || k == "page_size" || k == "fctime" || k == "tctime" {
				continue
			}
			tv := reflect.ValueOf(v)
			switch tv.Kind() {
			case reflect.Slice:
				slen := tv.Len()
				if slen == 1 {
					sql := fmt.Sprintf("%s=?", k)
					db = db.Where(sql, tv.Index(0).Interface())
				} else {
					sql := fmt.Sprintf("%s in ?", k)
					db = db.Where(sql, v)
				}
			default:
				sql := fmt.Sprintf("%s = ?", k)
				db = db.Where(sql, v)
			}
		}
		//时间范围
		var f_tm int64 = 0
		var t_tm int64 = 0
		if f_time, ok := boh.Qdata["fctime"]; ok {
			f_tm = (f_time.(time.Time)).Unix()
		}
		if t_time, ok := boh.Qdata["tctime"]; ok {
			t_tm = (t_time.(time.Time)).Unix()
		}
		logger.Debugf(ctx, "f_tm: %d t_tm: %d", f_tm, t_tm)
		if f_tm > 0 && t_tm > 0 {
			db = db.Where("ctime >= ? and ctime <= ?", f_tm, t_tm)
		}

		var count int64 = 0
		cdb := db
		cdb.Count(&count)
		logger.Debugf(ctx, "count: %d", count)

		it_page, _ := boh.Qdata["page"]
		it_page_size, _ := boh.Qdata["page_size"]
		page := it_page.(int)
		page_size := it_page_size.(int)

		var pages int64 = 0
		p_div := count % int64(page_size)
		if p_div == 0 {
			pages = count / int64(page_size)
		} else {
			pages = count/int64(page_size) + 1
		}
		logger.Debugf(ctx, "page: %d page_size: %d", page, page_size)
		db = db.Limit(page_size).Offset((page - 1) * page_size)
		//db = db.Order("ctime desc")
		ddatas := []map[string]interface{}{}
		ret := db.Scan(&ddatas)
		if ret.Error != nil {
			logger.Warnf(ctx, "query data: %s", ret.Error.Error())
			return nil, respcode.ERR_DB
		}
		re := boh.convertRows(ctx, ddatas)

		all := map[string]interface{}{
			"page":     page,
			"pagesize": page_size,
			"total":    count,
			"pagenum":  pages,
			"data":     re,
		}
		//return re, respcode.OK
		return all, respcode.OK
	}
	return nil, respcode.ERR_PARAM
}

func (boh *BaseObjectHandler) Post(ctx context.Context) error {
	var e_code int = respcode.OK
	data := map[string]interface{}{}
	q_str := boh.C.Param("base_edit")
	logger.Debugf(ctx, "q_str: %s", q_str)
	switch q_str {
	case "add":
		data, e_code = boh.Insert(ctx)
	case "mod":
		data, e_code = boh.Update(ctx)
	case "delete":
		logger.Debugf(ctx, "delete")
		data, e_code = boh.Delete(ctx)
	}
	if e_code != respcode.OK {
		return respcode.RetError[string](boh.C, e_code, "", "", "")
	}
	return respcode.RetSucc[map[string]interface{}](boh.C, data)
}

func (boh *BaseObjectHandler) Get(ctx context.Context) error {
	//return nil
	var e_code int = respcode.OK
	//ldata := []map[string]interface{}{}
	data := map[string]interface{}{}
	q_str := boh.C.Param("base_query")
	switch q_str {
	case "qlist":
		data, e_code = boh.QlistFunc(ctx)
	case "q":
		data, e_code = boh.GetData(ctx)

	}
	if e_code != respcode.OK {
		return respcode.RetError[string](boh.C, e_code, "", "", "")
	}
	/*if q_str == "qlist" {
		return respcode.RetSucc[[]map[string]interface{}](boh.C, ldata)
	}
	return respcode.RetSucc[map[string]interface{}](boh.C, data)*/
	return respcode.RetSucc[map[string]interface{}](boh.C, data)

}
