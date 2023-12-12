package user

import (
	"context"
	"usercenter/respcode"
	ut "usercenter/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/lanwenhong/lgobase/logger"
)

type TestQueryIds struct {
	Ids []int64 `form:"ids" binding:"required"`
}

func (uch *UserHandler) QueryIds(ctx context.Context) error {
	tqi := TestQueryIds{}
	if err := uch.C.ShouldBindWith(&tqi, binding.Query); err != nil {
		logger.Debugf(ctx, "get_user TestQueryIds data: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(tqi, err), "", "")
	}
	logger.Debugf(ctx, "ids: %v", tqi.Ids)
	ret := map[string]interface{}{}
	return respcode.RetSucc[map[string]interface{}](uch.C, ret)
}
