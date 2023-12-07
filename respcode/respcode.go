package respcode

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	OK           = 0
	ERR          = -1
	ERR_USER     = -2
	ERR_PARAM    = -3
	ERR_AUTH     = -4
	ERR_ACTION   = -5
	ERR_DATA     = -6
	ERR_PERM     = -7
	ERR_INTERNAL = -8
	ERR_DB       = -9
)

var ErrMap map[int]string = map[int]string{
	OK:           "成功",
	ERR:          "失败",
	ERR_USER:     "用户错误",
	ERR_PARAM:    "参数错误",
	ERR_AUTH:     "认证失败",
	ERR_ACTION:   "操作失败",
	ERR_DATA:     "数据错误",
	ERR_PERM:     "权限错误",
	ERR_INTERNAL: "内部错误",
	ERR_DB:       "数据操作错误",
}

const (
	// 未认证
	STATUS_NOAUTH = 1
	// 正常
	STATUS_OK = 2
	// 封禁
	STATUS_BAN = 3
	// 删除
	STATUS_DEL = 4
)

const (
	ROLE_PLATFORM_ADMIN = 1
	ROLE_ORG_ADMIN      = 2
	ROLE_ORG_USER       = 3
	ROLE_USER_NORMAL    = 4
)

var RoleMap map[int]string = map[int]string{
	ROLE_PLATFORM_ADMIN: "平台管理员",
	ROLE_ORG_ADMIN:      "机构管理员",
	ROLE_ORG_USER:       "机构用户",
	ROLE_USER_NORMAL:    "普通用户",
}

type UserCenterErr struct {
	Code int
	Msg  string
}

func NewUserCenterErr(code int, msg string) *UserCenterErr {
	err_msg := ""
	if msg == "" {
		err_msg = ErrMap[code]
	} else {
		err_msg = msg
	}
	return &UserCenterErr{
		Code: code,
		Msg:  err_msg,
	}
}

func (e *UserCenterErr) Error() string {
	return fmt.Sprintf("code: %d msg: %s", e.Code, e.Msg)
}

func RespError[T string | map[string]interface{}](errcode int, resperr string, respmsg string, data T) map[string]interface{} {
	resp := map[string]interface{}{
		"respcd":  errcode,
		"respmsg": respmsg,
		"resperr": resperr,
		"data":    data,
	}
	return resp

}

func RespSucc[T string | map[string]interface{}](resperr int, data T) map[string]interface{} {

	resp := map[string]interface{}{
		"respcd":  OK,
		"respmsg": "",
		"resperr": resperr,
		"data":    data,
	}
	return resp
}

func RetError[T string | map[string]interface{}](c *gin.Context, errcode int, resperr string, respmsg string, data T) error {
	err := NewUserCenterErr(errcode, resperr)
	resp := RespError[T](err.Code, err.Msg, respmsg, data)
	c.JSON(http.StatusOK, resp)
	return err
}

func RetSucc[T string | map[string]interface{}](c *gin.Context, data T) error {
	resp := RespSucc[T](OK, data)
	c.JSON(http.StatusOK, resp)
	return nil
}
