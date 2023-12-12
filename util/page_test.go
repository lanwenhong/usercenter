package util

import (
	"context"
	"testing"
	"time"

	"github.com/lanwenhong/lgobase/dbenc"
	"github.com/lanwenhong/lgobase/dbpool"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"

	//"github.com/lanwenhong/lgobase/util"
	ut "usercenter/util"

	dlog "gorm.io/gorm/logger"
)

func TestPageSplit(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenKsuid())
	myconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		ColorFull:    true,
		Loglevel:     logger.DEBUG,
	}

	logger.Newglog("./", "test.log", "test.log.err", myconf)
	dconfig := &dlog.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		LogLevel:                  dlog.Info,   // 日志级别
		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  true,        // 禁用彩色打印
	}

	db_conf := dbenc.DbConfNew(ctx, "../db.ini")
	dbs := dbpool.DbpoolNew(db_conf)
	dbs.SetormLog(ctx, dconfig)
	tk := "qfconf://usercenter?maxopen=1000&maxidle=30"
	err := dbs.Add(ctx, "usercenter", tk, dbpool.USE_GORM)
	if err != nil {
		t.Fatal(err)
	}
	//tdb := dbs.OrmPools["usercenter"]
	sql := "select * from users order by ctime desc"
	pdd := ut.NewPageDataDb(ctx, dbs, sql, "", "usercenter")
	pg := ut.NewPage(ctx, 1, 2, pdd)
	err = pg.Split(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pg.Pdd.Data)
	t.Log(len(pg.Pdd.Data))
	t.Log(pg.Count)
}
