package app

import (
	"fmt"
	"os"

	"github.com/Unknwon/goconfig"
	"github.com/go-xorm/core"
	_ "github.com/lunny/godbc"
	"github.com/lunny/log"
	"gopkg.in/macaron.v1"
)

const (
	CfgPath = "./app.conf"
)

var (
	AppCfg           *goconfig.ConfigFile
	AppLog           *log.Logger
	AppDB            *core.DB
	globalInsertChan chan bool

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

	globalInsertChan = make(chan bool)

	//init log
	AppLog = log.New(os.Stderr, "", log.Ldefault())
	w := log.NewFileWriter(log.FileOptions{
		ByType: log.ByDay,
		Dir:    "./logs",
	})
	AppLog.SetOutput(w)

	//read db config
	var err error
	AppCfg, err = goconfig.LoadConfigFile(CfgPath)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CfgPath, err))
	}

	AppName = AppCfg.MustValue("app", "app_name", "webapp")
	RunMode = AppCfg.MustValue("app", "run_mode", "dev")

	HttpPort = AppCfg.MustInt(RunMode, "http_port")
	DbDriver = AppCfg.MustValue(RunMode, "db_driver")
	DbDriverConnstr = AppCfg.MustValue(RunMode, "db_driver_connstr")
	DbUsername = AppCfg.MustValue(RunMode, "db_username")
	DbPassword = AppCfg.MustValue(RunMode, "db_password")
	DbServer = AppCfg.MustValue(RunMode, "db_server")
	DbDatebase = AppCfg.MustValue(RunMode, "db_datebase")
	DbPort = AppCfg.MustInt(RunMode, "db_port")

	//init db engine
	if AppDB == nil {
		connString := fmt.Sprintf(DbDriverConnstr, DbServer,
			DbUsername, DbPassword, DbPort, DbDatebase)

		AppLog.Info(connString)
		var err error
		AppDB, err = core.Open(DbDriver, connString)

		if err != nil {
			AppLog.Error(err)
		}
	}

	//set runmode
	if RunMode == "prod" {
		macaron.Env = macaron.PROD
	}
}
