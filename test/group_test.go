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

func TestGroupAdd(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)

	for i := 10; i < 20; i++ {
		name := fmt.Sprintf("test%d_group", i)
		info := fmt.Sprintf("my test %d group", i)
		req := map[string]string{
			"name":     name,
			"info":     info,
			"parentid": "0",
		}
		header := map[string]string{}
		header["Content-Type"] = "application/x-www-form-urlencoded"
		cookie := fmt.Sprintf("%s=%s", "sid", sk)
		header["Cookie"] = cookie
		c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
		v, _, err := c.Post(ctx, "uc/v1/group/add", 3000, req, header)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Log(string(v))
	}
}

func TestGroupQ(t *testing.T) {
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
		"id": "7140953806602309647",
	}

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/group/q", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))
}

func TestGroupQlist(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)

	/*req := map[string]string{
		"page":      "3",
		"page_size": "2",
	}*/

	/*names := []string{
		"test6_group",
		"test7_group",
	}
	req := map[string]interface{}{
		"page":      "1",
		"page_size": "2",
		"names":     names,
	}*/

	/*ids := []string{
		"7140953806602309643",
		"7140953806602309644",
	}*/
	/*info := []string{
		"ggggg test jujingyi",
		//"my test 4 group",
	}*/
	req := map[string]interface{}{
		"page":      "1",
		"page_size": "2",
		"fctime":    "2023-12-14 14:41:00",
		"tctime":    "2023-12-18 16:41:00",
		//"info":      info,
		//"id": ids,
		"id": "7140953806602309644",
	}

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	logger.Debugf(ctx, "cookie: %s", cookie)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Get(ctx, "uc/v1/group/qlist", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))

}

func TestGroupDel(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)
	/*ids := []string{
		"7140953806602309643",
		"7140953806602309644",
	}*/
	req := map[string]interface{}{
		"id": "7140953806602309648",
		//"ids": ids,
	}
	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	header["Cookie"] = cookie

	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Post(ctx, "uc/v1/group/delete", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))

}

func TestGroupUpdate(t *testing.T) {
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       true,
		Loglevel:     logger.DEBUG,
		ColorFull:    true,
	}
	logger.Newglog("./", "test.log", "test.log.err", lconf)

	ctx := context.Background()
	sk := getCookie(ctx)

	header := map[string]string{}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	cookie := fmt.Sprintf("%s=%s", "sid", sk)
	header["Cookie"] = cookie

	/*ids := []string{
		"7140953806602309641",
		"7140953806602309642",
	}*/
	req := map[string]interface{}{
		"id": "7140953806602309641",
		//"id":   ids,
		"info": "ggggg test jujingyi",
	}
	c := ghttpclient.QfHttpClientNew(ghttpclient.INTER_PROTO_PUSHAPI, "127.0.0.1:8000", false)
	v, _, err := c.Post(ctx, "uc/v1/group/mod", 3000, req, header)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(v))
}
