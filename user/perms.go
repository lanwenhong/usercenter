package user

import (
	"context"
	"time"
	"usercenter/respcode"
	ut "usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type PermsAddData struct {
	Name   string `form:"name" binding:"required" reg_error_info:"权限名字格式错误"`
	Info   string `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Detail string `form:"detail" binding:"omitempty" reg_error_info:"权限细节格式错误"`
}

type PermsDelData struct {
	Id uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
}

type PermsModDataWithId struct {
	Id   uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
	Info string `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name string `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
}

type PermsModDataWithIds struct {
	Id   []uint64 `form:"id" binding:"required" reg_error_info:"ids格式错误"`
	Info string   `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name string   `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
}

type PermsDataList struct {
	Id       []uint64  `form:"id" binding:"omitempty" reg_error_info:"ids格式错误"`
	Info     []string  `form:"info" binding:"omitempty" reg_error_info:"权限信息格式错误"`
	Name     []string  `form:"name" binding:"omitempty" reg_error_info:"权限名字格式错误"`
	Fctime   time.Time `form:"fctime" binding:"omitempty" time_format:"2006-01-02 15:04:05" reg_error_info:"ctime起始格>式错误"`
	Tctime   time.Time `form:"tctime" binding:"omitempty" time_format:"2006-01-02 15:04:05" reg_error_info:"ctime结束格>式错误"`
	Page     int       `form:"page" binding:"required" reg_error_info:"page格式错误"`
	PageSize int       `form:"page_size" binding:"required" reg_error_info:"page_size格式错误"`
}

type PermData struct {
	Id uint64 `form:"id" binding:"omitempty" reg_error_info:"id格式错误"`
}

type PermsOpHandler struct {
	BaseObjectHandler
	BaseOpFuncIndex map[string]BaseOpFunc
}

func (poh *PermsOpHandler) AddOpFunc(ctx context.Context) error {
	pad := PermsAddData{}
	if err := poh.C.ShouldBindWith(&pad, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pad, err), "", "")
	}
	poh.Qdata, _ = ut.Stru2Map(ctx, pad)
	return poh.Post(ctx)
}

func (poh *PermsOpHandler) DelOpFunc(ctx context.Context) error {
	pdd := PermsDelData{}
	if err := poh.C.ShouldBindWith(&pdd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pdd, err), "", "")
	}
	poh.Qdata, _ = ut.Stru2Map(ctx, pdd)
	logger.Debugf(ctx, "%v", poh.Qdata)
	return poh.Post(ctx)
}

func (poh *PermsOpHandler) ModOpFunc(ctx context.Context) error {
	pids := poh.C.PostFormArray("id")
	logger.Debugf(ctx, "pids: %v", pids)
	if len(pids) == 1 {
		pmd := PermsModDataWithId{}
		if err := poh.C.ShouldBindWith(&pmd, binding.Form); err != nil {
			logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
			return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pmd, err), "", "")
		}
		poh.Qdata, _ = ut.Stru2Map(ctx, pmd)
		return poh.Post(ctx)
	} else if len(pids) > 1 {
		pmd := PermsModDataWithIds{}
		if err := poh.C.ShouldBindWith(&pmd, binding.Form); err != nil {
			logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
			return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pmd, err), "", "")
		}
		poh.Qdata, _ = ut.Stru2Map(ctx, pmd)
		return poh.Post(ctx)
	}
	return respcode.RetError[string](poh.C, respcode.ERR, "modify data err", "", "")
}

func (poh *PermsOpHandler) QlistOpFunc(ctx context.Context) error {
	pdl := PermsDataList{}
	if err := poh.C.ShouldBindWith(&pdl, binding.Query); err != nil {
		logger.Warnf(ctx, "qlist binding data: %s", err.Error())
		return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pdl, err), "", "")
	}
	logger.Debugf(ctx, "qlist handler")
	poh.Qdata, _ = ut.Stru2Map(ctx, pdl)
	return poh.Get(ctx)
}

func (poh *PermsOpHandler) QopFunc(ctx context.Context) error {
	pd := PermData{}
	if err := poh.C.ShouldBindWith(&pd, binding.Query); err != nil {
		logger.Warnf(ctx, "q binding data: %s", err.Error())
		return respcode.RetError[string](poh.C, respcode.ERR, ut.ValidatErr(pd, err), "", "")
	}
	logger.Debugf(ctx, "q handler")
	poh.Qdata, _ = ut.Stru2Map(ctx, pd)
	return poh.Get(ctx)
}

func PermsOpHandlerNew(c *gin.Context, cookie string) *PermsOpHandler {
	poh := PermsOpHandler{}
	poh.C = c
	poh.Cookie = cookie
	poh.Table = "perms"

	poh.BaseOpFuncIndex = map[string]BaseOpFunc{
		"add":    poh.AddOpFunc,
		"mod":    poh.ModOpFunc,
		"delete": poh.DelOpFunc,
		"list":   poh.QlistOpFunc,
		"q":      poh.QopFunc,
	}
	poh.InsertFunc = poh.Insert
	poh.UpdateFunc = poh.Update
	poh.DeleteFunc = poh.Delete
	poh.QlistFunc = poh.GetDataList
	poh.Qfunc = poh.GetData
	return &poh
}

func PermsOp(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_edit")
	poh := PermsOpHandlerNew(c, cookie.(string))
	logger.Debugf(ctx, "grouped: %s", grouped)
	if op, ok := poh.BaseOpFuncIndex[grouped]; ok {
		op(ctx)
	} else {
		respcode.RetError[string](c, respcode.ERR, "not found method", "", "")
	}
}

func PermsQuery(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_query")
	logger.Debugf(ctx, "grouped: %s", grouped)
	poh := PermsOpHandlerNew(c, cookie.(string))
	if op, ok := poh.BaseOpFuncIndex[grouped]; ok {
		op(ctx)
	} else {
		respcode.RetError[string](c, respcode.ERR, "not found method", "", "")
	}
}
