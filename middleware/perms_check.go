package middleware

import (
	"context"
	"usercenter/respcode"
	"usercenter/session"
	ut "usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func CheckPerms(perms []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		ctx := context.WithValue(context.Background(), "trace_id", requestID)
		logger.Debugf(ctx, "check perms middleware")
		logger.Debugf(ctx, "check perms middleware perms: %v", perms)
		sk, _ := c.Get("have_se")
		logger.Debugf(ctx, "sk: %s", sk)
		se := session.NewSession(sk.(string))
		value, _ := se.GetData(ctx, "allperm")
		in_allperm := value.([]interface{})
		logger.Debugf(ctx, "allperm: %v len: %d", in_allperm, len(in_allperm))
		in_isadmin, _ := se.GetData(ctx, "isadmin")
		isadmin := int64(in_isadmin.(float64))
		logger.Debugf(ctx, "isadmin: %d", isadmin)
		allperm := []string{}
		for _, v := range in_allperm {
			allperm = append(allperm, v.(string))
		}
		user_perm_set := ut.NewSet[string]()
		user_perm_set.SetList(ctx, allperm)
		check_perm_set := ut.NewSet[string]()
		check_perm_set.SetList(ctx, perms)

		if is_sub := user_perm_set.IsSubSet(ctx, check_perm_set); !is_sub || isadmin != 1 {
			respcode.RetError[string](c, respcode.ERR, "permission deny", "", "")
			c.Abort()
		} else {

			c.Next()
		}
	}
}
