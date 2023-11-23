package sruntime

import (
	"context"
	"time"
	"usercenter/myconf"

	"github.com/lanwenhong/lgobase/dbenc"
	"github.com/lanwenhong/lgobase/dbpool"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/redispool"
	"github.com/redis/go-redis/v9"
)

type SrunTime[T redispool.RedisMethod] struct {
	Dbs *dbpool.Dbpool
	Rop *redispool.RedisOP[T]
}

var Gsvr *SrunTime[*redis.Client]

func CreateRuntime(ctx context.Context) (*SrunTime[*redis.Client], error) {
	g_svr := new(SrunTime[*redis.Client])
	logger.Debugf(ctx, "tk: %s", myconf.Scnf.TokenFile)
	db_conf := dbenc.DbConfNew(ctx, myconf.Scnf.TokenFile)
	g_svr.Dbs = dbpool.DbpoolNew(db_conf)
	g_svr.Dbs.Add(ctx, "usercenter", myconf.Scnf.UcToken, dbpool.USE_GORM)
	logger.Debugf(ctx, "redis db: %d addr: %s pool_size: %d minidle: %d", myconf.Scnf.Db, myconf.Scnf.RedisAddr, myconf.Scnf.PoolSize, myconf.Scnf.MinIdle)

	rdb := redispool.NewGrPool(ctx, "", "", myconf.Scnf.Db, myconf.Scnf.RedisAddr, myconf.Scnf.PoolSize, myconf.Scnf.MinIdle,
		time.Duration(myconf.Scnf.ConnTimeout)*time.Second,
		time.Duration(myconf.Scnf.ReadTimeout)*time.Second,
		time.Duration(myconf.Scnf.WriteTimeout)*time.Second,
	)
	g_svr.Rop = redispool.NewRedisOP[*redis.Client](rdb)
	//op := redispool.NewRedisOP[*redis.Client](rdb)
	//op.Rdb.Ping(ctx).Result()

	Gsvr = g_svr
	return g_svr, nil
}
