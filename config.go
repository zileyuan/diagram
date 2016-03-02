package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/lunny/log"
	"github.com/go-macaron/macaron"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const (
	CfgPath = "./app.conf"
)

var (
	Cfg *goconfig.ConfigFile
        Log *log.Logger
	Orm *xorm.Engine

	AppName string
	RunMode string

	HttpPort        int
	DbDriver        string
	DbDriverConnstr string
	DbUsername      string
	DbPassword      string
	DbServer        string
	DbPort          int
	DbDatebase      string
)

func init() {

	//init log
	Log = log.New(os.Stderr, "", log.Ldefault())
	w := log.NewFileWriter(log.FileOptions{
		ByType: log.ByDay,
		Dir:    "./logs",
	})
	Log.SetOutput(w)

	//read db config
	var err error
	Cfg, err = goconfig.LoadConfigFile(CfgPath)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CfgPath, err))
	}

	AppName = Cfg.MustValue("app", "app_name", "webapp")
	RunMode = Cfg.MustValue("app", "run_mode", "dev")

	HttpPort = Cfg.MustInt(RunMode, "http_port")
	DbDriver = Cfg.MustValue(RunMode, "db_driver")
	DbDriverConnstr = Cfg.MustValue(RunMode, "db_driver_connstr")
	DbUsername = Cfg.MustValue(RunMode, "db_username")
	DbPassword = Cfg.MustValue(RunMode, "db_password")
	DbServer = Cfg.MustValue(RunMode, "db_server")
	DbDatebase = Cfg.MustValue(RunMode, "db_datebase")
	DbPort = Cfg.MustInt(RunMode, "db_port")

	//init db engine
	if Orm == nil {
		connString := fmt.Sprintf(DbDriverConnstr, DbServer,
			DbPort, DbUsername, DbPassword, DbDatebase)

		Log.Info(connString)
		var err error
		Orm, err = xorm.NewEngine(DbDriver, connString)

		if err == nil {
			Orm.TZLocation = time.Local
			Orm.ShowSQL(true)
			Orm.SetMapper(core.SameMapper{})
		} else {
			Log.Error(err)
		}
	}

	//set runmode
	if RunMode == "prod" {
		macaron.Env = macaron.PROD
	}
}
