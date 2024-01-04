package router

import (
	"usercenter/middleware"
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
	perms := []string{
		"perm_view",
	}
	r.GET("uc/v1/perm/:base_query", middleware.CheckPerms(perms), user.PermsQuery)
	perms = []string{
		"perm_mod",
	}
	r.POST("uc/v1/perm/:base_edit", middleware.CheckPerms(perms), user.PermsOp)

	r.GET("/uc/v1/role/:base_query", user.RoleQuery)
	r.POST("/uc/v1/role/:base_edit", user.RoleOp)
}
