package middleware

import (
	"context"
	"usercenter/session"

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
		allperm := value.([]interface{})
		logger.Debugf(ctx, "allperm: %v", allperm)
		c.Next()
	}
}
