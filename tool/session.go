package tool

import (
	"context"
	"net/http"
	"usercenter/respcode"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
)

func SessionTest(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	ctx := context.WithValue(context.Background(), "trace_id", requestID)

	logger.Debugf(ctx, "set cookie")
	c.SetCookie("sid", "Ilovenannan", 3600, "/", "127.0.0.1:8000",
		false, false)
	resp := respcode.RespSucc[string](respcode.OK, "")
	c.JSON(http.StatusOK, resp)

	return
}
