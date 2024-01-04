package usercenter

import (
	"context"
	"fmt"
	"testing"

	"github.com/lanwenhong/lgobase/ghttpclient"
	"github.com/lanwenhong/lgobase/logger"
)

func getCookie(ctx context.Context) string {
	req := map[string]interface{}{
		"mobile":   "13012340000",
		"password": "111111",
	}
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	_, resp, err := c.Get(ctx, "uc/v1/user/login", 3000, req, header)
	if err != nil {
		logger.Warnf(ctx, err.Error())
		return ""
	}
	sk := ""
	for _, c := range resp.Cookies() {
		sk = c.Value
	}
	return sk
}

func TestPermsAdd(t *testing.T) {

	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("test%d_perm", i)
		info := fmt.Sprintf("my test %d perm", i)
		detail := fmt.Sprintf("perm_%d jujingyi", i)
		req := map[string]string{
			"name":   name,
			"info":   info,
			"detail": detail,
		}
		header := map[string]string{}
		header["Content-Type"] = "application/x-www-form-urlencoded"
		cookie := fmt.Sprintf("%s=%s", "sid", sk)
		header["Cookie"] = cookie
		c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
		v, _, err := c.Post(ctx, "uc/v1/perm/add", 3000, req, header)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Log(string(v))
	}
}

func TestPermQ(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)

	req := map[string]interface{}{
		"id": "1",
	}

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/perm/q", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))
}
