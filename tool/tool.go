package tool

import (
	"context"
	"image/color"
	"net/http"
	"time"
	"usercenter/dbmodel"
	"usercenter/respcode"
	"usercenter/sruntime"

	"github.com/gin-gonic/gin"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/mojocn/base64Captcha"
)

type CodeVerifyForm struct {
	Id     string `form:"id" binding:"required,len=20"`
	Answer string `form:"answer" binding:"required,len=4"`
}

type DbStore struct {
	traceid string
}

type ImgMake struct {
}

func (r *DbStore) Set(vid string, value string) error {
	ctx := context.WithValue(context.Background(), "trace_id", r.traceid)
	id, _ := util.Genid(ctx, sruntime.Gsvr.Dbs.OrmPools["usercenter"])
	vcode := dbmodel.VerifyCode{
		Id:       id,
		VerifyId: vid,
		Answer:   value,
		Valid:    1,
		Stime:    uint64(time.Now().Unix()),
		Etime:    uint64(time.Now().Unix() + 120),
	}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.WithContext(ctx).Create(&vcode)
	return ret.Error
}

func (s *DbStore) Verify(id, answer string, clear bool) bool {
	//return true
	v := s.Get(id, true)

	if v == "" {
		return false
	}
	return v == answer
}

func (s *DbStore) clear(code string) error {
	ctx := context.WithValue(context.Background(), "trace_id", s.traceid)
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.Where("verify_id = ?", code).Delete(dbmodel.VerifyCode{})
	if ret.Error != nil {
		logger.Warnf(ctx, "delete %s err %s", code, ret.Error.Error())
	}
	return ret.Error
}

func (s *DbStore) Get(id string, clear bool) (value string) {
	if clear {
		defer s.clear(id)
	}
	ctx := context.WithValue(context.Background(), "trace_id", s.traceid)
	vcode := dbmodel.VerifyCode{}
	db := sruntime.Gsvr.Dbs.OrmPools["usercenter"]
	ret := db.WithContext(ctx).Select("answer", "stime", "etime").Where("verify_id=? and valid=?", id, 1).Find(&vcode)
	if ret.Error != nil {
		logger.Warnf(ctx, "get answer err: %s", ret.Error.Error())
		return ""
	}
	now := uint64(time.Now().Unix())
	logger.Debugf(ctx, "now: %d stime: %d etime: %d", now, vcode.Stime, vcode.Etime)

	if vcode.Stime <= now && now <= vcode.Etime {
		return vcode.Answer
	}
	logger.Debug(ctx, "verify code expired")
	return ""
}

func (im ImgMake) MakeOne(c *gin.Context) (vid, b64s string, err error) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	requestID := c.Request.Header.Get("X-Request-ID")
	store := DbStore{
		traceid: requestID,
	}
	captcha := base64Captcha.NewCaptcha(driver, &store)
	lid, lb64s, lerr := captcha.Generate()
	return lid, lb64s, lerr
}

func GetImageCode(c *gin.Context) {
	im := ImgMake{}
	lid, lb64s, lerr := im.MakeOne(c)
	if lerr != nil {
		resp := respcode.RespError[string](respcode.ERR, lerr.Error(), "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	data := map[string]interface{}{
		"id":   lid,
		"b64s": lb64s,
	}
	resp := respcode.RespSucc[map[string]interface{}](respcode.OK, data)
	c.JSON(http.StatusOK, resp)
}

func CodeVerify(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	//ctx := context.WithValue(context.Background(), "trace_id", requestID)

	cf := CodeVerifyForm{}
	if err := c.ShouldBindQuery(&cf); err != nil {
		resp := respcode.RespError[string](respcode.ERR, err.Error(), "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	store := DbStore{
		traceid: requestID,
	}
	ret := store.Verify(cf.Id, cf.Answer, true)
	if !ret {
		resp := respcode.RespError[string](respcode.ERR, "验证码不正确，请重新输入", "", "")
		c.JSON(http.StatusOK, resp)
		return
	}
	resp := respcode.RespSucc[string](respcode.OK, "验证码校验成功")
	c.JSON(http.StatusOK, resp)
	return
}
