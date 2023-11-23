package util

import (
	"testing"

	"usercenter/util"

	"github.com/lanwenhong/lgobase/logger"
)

func TestEmailValidator(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	//ctx := context.Background()
	ret := util.VerifyEmailFormat("hexiexuanlv@126.com")
	if ret {
		t.Log("succ")
	} else {
		t.Log("fail")
	}

}
