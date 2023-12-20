package util

import (
	"context"
	"testing"
	"usercenter/util"

	"github.com/lanwenhong/lgobase/logger"
)

func TestSetList(t *testing.T) {

	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)
	ctx := context.Background()
	set := util.NewSet[string]()
	ss := []string{
		"jujingyi",
		"chenduling",
		"liushishi",
	}
	set.SetList(ctx, ss)
	logger.Debugf(ctx, "%v", set.InMap)

	set2 := util.NewSet[string]()
	ss2 := []string{
		"jujingyi",
		"chenduling",
		"liushishi",
		"zhaoliying",
	}
	set2.SetList(ctx, ss2)
	logger.Debugf(ctx, "%v", set2.InMap)
	ret := set2.IsSubSet(ctx, set)
	t.Log(ret)

}
