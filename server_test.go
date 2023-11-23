package usercenter

import (
	"context"
	"testing"

	"github.com/lanwenhong/lgobase/ghttpclient"
	"github.com/lanwenhong/lgobase/logger"
)

func TestUserRegister(t *testing.T) {

	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	req := map[string]string{
		"username": "zhao",
		"mobile":   "18010583873",
		"email":    "hexiexuanlv@126.com",
		"password": "111111",
	}

	ctx := context.Background()
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	//header["Cookie"] = "sid=xxx"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, resp, err := c.Post(ctx, "uc/v1/user/signup", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	for k, v := range resp.Cookies() {
		t.Log(k)
		t.Log(v.Name)
		t.Log(v.Value)
		t.Log(v.Path)
		t.Log(v.Expires)
	}
	t.Log(string(v))

}

func TestUserlogin(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	req := map[string]string{
		//"username": "zhao",
		//"mobile": "13012340000",
		"email":    "hexiexuanlv1@126.com",
		"password": "111111",
	}

	ctx := context.Background()
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	//header["Cookie"] = "sid=xxx"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, resp, err := c.Get(ctx, "uc/v1/user/login", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, c := range resp.Cookies() {
		t.Log(c.Name)
		t.Log(c.Value)
		t.Log(c.Path)
		t.Log(c.Expires)
	}
	t.Log(string(v))
}
