package session

import (
	"context"
	"encoding/json"
	"time"
	"usercenter/sruntime"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

type SeHandler interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, map[string]interface{}, time.Duration) error
}

type Session struct {
	Key string
}

func NewSession(key string) *Session {
	return &Session{
		Key: key,
	}
}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	ret, err := sruntime.Gsvr.Rop.Rdb.Get(ctx, key).Result()

	if err != nil {
		logger.Warnf(ctx, "get err: %s", err.Error())
		return "", err
	}
	return ret, nil
}

func (s *Session) Set(ctx context.Context, key string, value map[string]interface{}, expire time.Duration) error {
	v, _ := json.Marshal(value)
	_, err := sruntime.Gsvr.Rop.Rdb.SetEx(ctx, key, v, expire*time.Second).Result()

	if err != nil {
		logger.Warnf(ctx, "set err: %s", err.Error())
		return err
	}
	return nil
}

func (s *Session) Update(ctx context.Context, key string, value map[string]interface{}, expire time.Duration) (string, error) {
	sek := ""
	if key == "" {
		sek = util.GenKsuid()
	} else {
		sek = key
	}
	v, _ := json.Marshal(value)
	_, err := sruntime.Gsvr.Rop.Rdb.SetEx(ctx, sek, v, expire*time.Second).Result()
	if err != nil {
		logger.Warnf(ctx, "set err: %s", err.Error())
		return "", err
	}
	return sek, nil
}
