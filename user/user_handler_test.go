package user

import (
	"context"
	"testing"
	"usercenter/myconf"
	"usercenter/sruntime"
	"usercenter/user"

	"github.com/lanwenhong/lgobase/logger"
)

func TestGetUser(t *testing.T) {
	conf_file := "/home/lanwenhong/dev/usercenter/usercenter.ini"
	myconf.Parseconf(conf_file)

	loglevel, _ := logger.LoggerLevelIndex(myconf.Scnf.LogLevel)
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       myconf.Scnf.LogStdOut,
		Loglevel:     loglevel,
		ColorFull:    myconf.Scnf.Colorfull,
	}
	logger.Newglog(myconf.Scnf.LogDir, myconf.Scnf.LogFile, myconf.Scnf.LogFileErr, lconf)

	ctx := context.Background()
	sruntime.CreateRuntime(ctx)

	uh := user.UserHandler{}

	ret, err := uh.TestgetUser(ctx, 7060864841283604340)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
