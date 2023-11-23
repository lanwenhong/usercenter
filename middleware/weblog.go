package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

func LoogerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		handleStart := time.Now()
		requestID := util.GenKsuid()
		c.Request.Header.Set("X-Request-ID", requestID)
		ctx := context.WithValue(context.Background(), "trace_id", requestID)

		var reqData []byte
		if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPost {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			logger.Debugf(ctx, "bodyBytes: %s", string(bodyBytes))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) >= 1024 {
				reqData = bodyBytes[0:1024]
			} else {
				reqData = bodyBytes
			}
		}
		// 处理请求
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		// 请求路由
		//reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()

		path := c.Request.URL.Path

		rawQuery := c.Request.URL.RawQuery

		if rawQuery == "" {
			rawQuery = "-"
		}
		//ctx := context.Background()
		//q_str := c.Request.URL.Query()
		//logger.Debugf(ctx, "q_str: %s", q_str)
		logger.Debugf(ctx, "path: %s", path)
		logger.Debugf(ctx, "query: %s", rawQuery)

		// 请求IP
		clientIp := c.ClientIP()
		handleEnd := time.Now()
		handleTime := handleEnd.Sub(handleStart)
		//logger.Infof(ctx, "|%3d|%13v|%15s|%s|%s|",
		logger.Infof(ctx, "%d|%v|%v|%s|%s|%s|%s|%s",
			statusCode,
			latencyTime,
			handleTime,
			clientIp,
			reqMethod,
			path,
			rawQuery,
			reqData,
		)
	}
}
