package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"usercenter/middleware"
	"usercenter/myconf"
	"usercenter/router"
	"usercenter/sruntime"
	"usercenter/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lanwenhong/lgobase/logger"
	//"github.com/redis/go-redis/v9"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	arg_num := len(os.Args)
	if arg_num != 2 {
		fmt.Printf("input param error")
		return
	}
	var filename = os.Args[1]
	err := myconf.Parseconf(filename)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	loglevel, _ := logger.LoggerLevelIndex(myconf.Scnf.LogLevel)
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       myconf.Scnf.LogStdOut,
		Loglevel:     loglevel,
		ColorFull:    myconf.Scnf.Colorfull,
	}
	logger.Newglog(myconf.Scnf.LogDir, myconf.Scnf.LogFile, myconf.Scnf.LogFileErr, lconf)

	ctx := context.Background()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("EmailValidator", util.EmailValidator)
		if err != nil {
			logger.Warnf(ctx, "register fail")
		}
	}
	logger.Debugf(ctx, "register succ")
	sruntime.CreateRuntime(ctx)
	logger.Debug(ctx, "server run")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(middleware.LoogerToFile())
	r.Use(middleware.CheckSession())
	router.Router(r)
	BindAddr := fmt.Sprintf("%s:%d", myconf.Scnf.Addr, myconf.Scnf.Port)
	r.Run(BindAddr)
}
