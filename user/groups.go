package user

import (
	"context"
	"net/http"
	"usercenter/respcode"
	ut "usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type GroupsAddData struct {
	Name     string `form:"name" binding:"omitempty" reg_error_info:"组名字格式错误"`
	Info     string `form:"info" binding:"omitempty" reg_error_info:"组信息格式错误"`
	parentid uint64 `form:"parentid" binding:"omitempty" reg_error_info:"父id格式错误"`
}

type GroupsDelData struct {
	Id uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
}

type GroupsModData struct {
	Id       uint64 `form:"id" binding:"required" reg_error_info:"id格式错误"`
	Info     string `form:"info" binding:"omitempty" reg_error_info:"组信息格式错误"`
	Name     string `form:"name" binding:"omitempty" reg_error_info:"组名字格式错误"`
	parentid uint64 `form:"parentid" binding:"omitempty" reg_error_info:"父id格式错误"`
}

type GroupsModDataList struct {
	Id        uint64   `form:"id" binding:"omitempty" reg_error_info:"id格式错误"`
	Ids       []uint64 `form:"ids" binding:"omitempty" reg_error_info:"ids格式错误"`
	Info      string   `form:"info" binding:"omitempty" reg_error_info:"组信息格式错误"`
	Name      string   `form:"name" binding:"omitempty" reg_error_info:"组名字格式错误"`
	Parentid  uint64   `form:"parentid" binding:"omitempty" reg_error_info:"父id格式错误"`
	Parentids []uint64 `form:"parentids" binding:"omitempty" reg_error_info:"Parentids格式错误"`
	Page      int      `form:"page" binding:"required" reg_error_info:"page格式错误"`
	PageSize  int      `form:"page_size" binding:"required" reg_error_info:"page_size格式错误"`
}

type GroupsOpHandler struct {
	BaseObjectHandler
	BaseOpFuncIndex map[string]BaseOpFunc
}

func (goh *GroupsOpHandler) AddOpFunc(ctx context.Context) error {
	gad := GroupsAddData{}
	if err := goh.C.ShouldBindWith(&gad, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gad, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gad)
	return goh.Post(ctx)
}

func (goh *GroupsOpHandler) DelOpFunc(ctx context.Context) error {
	gdd := GroupsDelData{}
	if err := goh.C.ShouldBindWith(&gdd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gdd, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gdd)
	return goh.Post(ctx)
}

func (goh *GroupsOpHandler) ModOpFunc(ctx context.Context) error {
	gmd := GroupsModData{}
	if err := goh.C.ShouldBindWith(&gmd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gmd, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gmd)
	return goh.Post(ctx)
}

func (goh *GroupsOpHandler) QlistOpFunc(ctx context.Context) error {
	gmpdl := GroupsModDataList{}
	if err := goh.C.ShouldBindWith(&gmpdl, binding.Query); err != nil {
		logger.Warnf(ctx, "qlist binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gmpdl, err), "", "")
	}
	logger.Debugf(ctx, "qlist handler")
	goh.Qdata, _ = ut.Stru2Map(ctx, gmpdl)
	return goh.Get(ctx)
}

func GroupsOpHandlerNew(c *gin.Context, cookie string) *GroupsOpHandler {
	gph := GroupsOpHandler{}
	gph.C = c
	gph.Cookie = cookie
	//gph.Table = "groups"
	gph.Table = "users"

	gph.BaseOpFuncIndex = map[string]BaseOpFunc{
		"add":    gph.AddOpFunc,
		"mod":    gph.ModOpFunc,
		"delete": gph.DelOpFunc,
		"qlist":  gph.QlistOpFunc,
	}
	return &gph
}

func GroupsOp(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_edit")
	se_check, _ := c.Get("check_session")
	if se_check.(string) == "fail" {
		resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	goh := GroupsOpHandlerNew(c, cookie.(string))
	op, _ := goh.BaseOpFuncIndex[grouped]
	op(ctx)
}

func GroupsQuery(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("base_query")
	logger.Debugf(ctx, "grouped: %s", grouped)
	se_check, _ := c.Get("check_session")
	if se_check.(string) == "fail" {
		resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	goh := GroupsOpHandlerNew(c, cookie.(string))
	op, _ := goh.BaseOpFuncIndex[grouped]
	op(ctx)
}
