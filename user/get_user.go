package user

import (
	"context"
	"usercenter/respcode"
	ut "usercenter/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type GetUserData struct {
	Userid uint64 `form:"userid" binding:"required" reg_error_info:"userid格式错误"`
}

func (uch *UserHandler) GetUser(ctx context.Context) error {
	gud := GetUserData{}
	if err := uch.C.ShouldBindWith(&gud, binding.Query); err != nil {
		logger.Debugf(ctx, "get_user binding data: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(mud, err), "", "")
	}
	user, err := uch.getUserBase(ctx, gud.Userid)
	if err != nil {
		logger.Warnf(ctx, "%s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, "get userinfo error", "", "")
	}
	return respcode.RetSucc[map[string]interface{}](uch.C, user)
}
