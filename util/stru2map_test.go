package util

import (
	"context"
	"testing"
	"usercenter/util"

	"github.com/lanwenhong/lgobase/logger"
)

type TestData struct {
	Id   int64  `form:"id"`
	Name string `form:"name"`
}

func TestEmailValidator(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()

	td := TestData{
		Id: 111,
		//Name: "ffffff",
	}
	ret, err := util.Stru2Map(ctx, td)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}
