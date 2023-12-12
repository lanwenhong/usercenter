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

type GroupsOpHandler struct {
	BaseObjectHandler
	BaseOpFuncIndex map[string]BaseOpFunc
}

func (goh *GroupsOpHandler) AddOpFunc(ctx context.Context) error {
	gad := GroupsAddData{}
	if err := uch.C.ShouldBindWith(&gad, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gad, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gad)
	return goh.Post(ctx)
}

func (goh *GroupsOpHandler) DelOpFunc(ctx context.Context) error {
	gdd := GroupsDelData{}
	if err := uch.C.ShouldBindWith(&gdd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gdd, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gad)
	return goh.Post(ctx)
}

func (goh *GroupsOpHandler) ModOpFunc(ctx context.Context) error {
	gmd := GroupsModData{}
	if err := uch.C.ShouldBindWith(&gmd, binding.Form); err != nil {
		logger.Warnf(ctx, "modify_user binding data: %s", err.Error())
		return respcode.RetError[string](goh.C, respcode.ERR, ut.ValidatErr(gmd, err), "", "")
	}
	goh.Qdata, _ = ut.Stru2Map(ctx, gad)
	return goh.Post(ctx)
}

func GroupsOpHandlerNew(c *gin.Context, cookie string) *GroupsOpHandler {
	gph := GroupsOpHandler{
		C:      c,
		Cookie: cookie.(string),
		Groups: "groups",
	}
	//gph.Qdata = map[string]interface{}{}
	gph.BaseOpFuncIndex = map[string]BaseOpFunc{
		"add":    gph.AddOpFunc,
		"mod":    gph.ModOpFunc,
		"delete": gph.DelOpFunc,
	}
}

func GroupsOp(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)
	cookie, _ := c.Get("have_se")
	grouped := c.Param("group_edit")
	se_check, _ := c.Get("check_session")
	if se_check.(string) == "fail" {
		resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	goh := GroupsOpHandlerNew(c, cookie)
	op, _ := gop.GroupsOpHandler[grouped]
	op(ctx)
}
