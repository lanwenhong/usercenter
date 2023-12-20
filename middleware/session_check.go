package middleware

import (
	"context"
	"usercenter/respcode"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func FindPath(paths []string, path string) bool {
	for _, v := range paths {
		if v == path {
			return true
		}
	}
	return false
}

func CheckSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		no_need := []string{
			"/uc/v1/user/signup",
			"/uc/v1/user/signup3rd",
			"/uc/v1/user/login",
			"/uc/v1/user/login3rd",
			"/uc/v1/user/login_reg_3rd",
			"/uc/v1/user/get_user",
			"/uc/v1/user/signin",
			"/uc/v1/user/query_ids",
			//"/uc/v1/group/qlist",
		}
		//c.Set("check_session", "succ")
		requestID := c.Request.Header.Get("X-Request-ID")
		ctx := context.WithValue(context.Background(), "trace_id", requestID)
		client_cookie, err := c.Cookie("sid")
		if err != nil {
			logger.Debugf(ctx, "not have cookie")
			//c.Request.Header.Set("X-have-se", "")
			c.Set("have_se", "")
			if !FindPath(no_need, c.Request.URL.Path) {
				logger.Debugf(ctx, "check session err")
				//c.Set("check_session", "fail")
				respcode.RetError[string](c, respcode.ERR, "session check error", "", "")
				c.Abort()
				return
				//resp := respcode.RespError[string](respcode.ERR, "session check error", "", "")
				//c.JSON(http.StatusOK, resp)
				//return
			} else {
				logger.Debugf(ctx, "no need check sid")
			}
		} else {
			logger.Debugf(ctx, "found client cookie: %s", client_cookie)
			c.Set("have_se", client_cookie)
		}
		c.Next()
	}
}
