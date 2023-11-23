package util

import (
	"testing"
	"usercenter/util"
)

func TestCreatePass(t *testing.T) {
	password := "I love zhangruonan"
	ret, err := util.CreatePassword(password)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(ret)
}
