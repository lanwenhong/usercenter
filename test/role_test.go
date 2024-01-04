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

func TestRoleList(t *testing.T) {
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
		"page":      "1",
		"page_size": "2",
		//"fctime":    "2023-12-14 14:41:00",
		//"tctime":    "2023-12-18 16:41:00",
		//"info":      info,
		//"id": ids,
		//"id": "7140953806602309644",
	}

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	logger.Debugf(ctx, "cookie: %s", cookie)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/role/list", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))

}

func TestRoleQ(t *testing.T) {
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
		"id": "7060862410235645810",
	}

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	logger.Debugf(ctx, "cookie: %s", cookie)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/role/q", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))

}
