package myconf

import (
	"fmt"

	"github.com/lanwenhong/lgobase/confparse"
	"github.com/lanwenhong/lgobase/gconfig"
)

type Connmap map[string]string

type Mcnf struct {
	Addr string `confpos:"server:addr" dtype:"base"`
	Port int    `confpos:"server:port" dtype:"base"`

	LogFile    string `confpos:"log:logfile" dtype:"base"`
	LogFileErr string `confpos:"log:logfile_err" dtype:"base"`
	LogDir     string `confpos:"log:logdir" dtype:"base"`
	LogLevel   string `confpos:"log:loglevel" dtype:"base"`
	LogStdOut  bool   `confpos:"log:logstdout" dtype:"base"`
	Colorfull  bool   `confpos:"log:colorfull" dtype:"base"`

	RedisAddr    string `confpos:"redis:redis_addr" dtype:"base"`
	PoolSize     int    `confpos:"redis:pool_size" dtype:"base"`
	MinIdle      int    `confpos:"redis:min_idle" dtype:"base"`
	ReadTimeout  int    `confpos:"redis:read_timeout" dtype:"base"`
	WriteTimeout int    `confpos:"redis:write_timeout" dtype:"base"`
	ConnTimeout  int    `confpos:"redis:connect_timeout" dtype:"base"`
	Db           int    `confpos:"redis:db" dtype:"base"`
	User         string `confpos:"redis:user" dtype:"base"`
	Passwd       string `confpos:"redis:passwd" dtype:"base"`

	UcToken   string `confpos:"db:usercenter" dtype:"base"`
	TokenFile string `confpos:"db:token_file" dtype:"base"`
}

var Scnf *Mcnf = new(Mcnf)

func Parseconf(filename string) error {
	cfg := gconfig.NewGconf(filename)
	err := cfg.GconfParse()
	if err != nil {
		fmt.Printf("parse %s %s", filename, err.Error())
		return err
	}
	cp := confparse.CpaseNew(filename)
	err = cp.CparseGo(Scnf, cfg)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("PushAddr: %s\n", Scnf.Addr)
	fmt.Printf("PushPort: %d\n", Scnf.Port)

	return err
}
