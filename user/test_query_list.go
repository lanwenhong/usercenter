package user

import (
	"context"
	"usercenter/respcode"

	"github.com/lanwenhong/lgobase/logger"
)

/*type TestQueryIds struct {
	Ids []int64 `form:"ids" binding:"required"`
}*/

func (uch *UserHandler) QueryIds(ctx context.Context) error {
	/*tqi := TestQueryIds{}
	if err := uch.C.ShouldBindWith(&tqi, binding.Query); err != nil {
		logger.Debugf(ctx, "get_user TestQueryIds data: %s", err.Error())
		return respcode.RetError[string](uch.C, respcode.ERR, ut.ValidatErr(tqi, err), "", "")
	}
	logger.Debugf(ctx, "ids: %v", tqi.Ids)*/
	x := uch.C.QueryArray("ids")
	logger.Debugf(ctx, "x: %v", x)
	queryParams := uch.C.Request.URL.Query()
	for k, v := range queryParams {
		logger.Debugf(ctx, "k: %s v: %s", k, v)
	}
	ret := map[string]interface{}{}
	return respcode.RetSucc[map[string]interface{}](uch.C, ret)
}
