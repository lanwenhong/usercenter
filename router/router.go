package router

import (
	"usercenter/tool"
	"usercenter/user"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/uc/v1/code/image", tool.GetImageCode)
	r.POST("/uc/v1/code/image", tool.GetImageCode)
	r.GET("/uc/v1/code/verify", tool.CodeVerify)
	r.POST("/uc/v1/code/verify", tool.CodeVerify)
	r.POST("/uc/v1/user/:useredit", user.UserOp)
	r.GET("/uc/v1/user/:userquery", user.UserQuery)
	r.POST("uc/v1/group/:base_edit", user.GroupsOp)
	r.GET("uc/v1/group/:base_query", user.GroupsQuery)
}
