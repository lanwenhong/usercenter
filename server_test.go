package usercenter

import (
	"context"
	"fmt"
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
		"mobile": "1301234000",
		//"email":    "hexiexuanlv1@126.com",
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

func TestGetUser(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	req := map[string]string{
		"userid": "7060864841283604340",
	}
	ctx := context.Background()
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/user/get_user", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))
}

func TestModifyUser(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	req := map[string]string{
		"mobile":   "13012340000",
		"password": "111111",
	}
	ctx := context.Background()
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, resp, err := c.Get(ctx, "uc/v1/user/login", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	sk := ""
	for _, c := range resp.Cookies() {
		t.Log(c.Name)
		t.Log(c.Value)
		t.Log(c.Path)
		t.Log(c.Expires)
		sk = c.Value
	}
	t.Log(string(v))

	mreq := map[string]string{
		"userid": "7133368332320837558",
		"mobile": "13800990890",
	}
	header = map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	header["Cookie"] = cookie
	c = ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, resp, err = c.Post(ctx, "uc/v1/user/mod", 3000, mreq, header)
	//v, resp, err = c.Get(ctx, "uc/v1/user/mod", 3000, mreq, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))
}

func TestAddGroup(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	req := map[string]string{
		"mobile":   "13012340000",
		"password": "111111",
	}
	ctx := context.Background()
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, resp, err := c.Get(ctx, "uc/v1/user/login", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	sk := ""
	for _, c := range resp.Cookies() {
		t.Log(c.Name)
		t.Log(c.Value)
		t.Log(c.Path)
		t.Log(c.Expires)
		sk = c.Value
	}
	t.Log(string(v))

}
