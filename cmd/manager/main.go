package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/config"
	"github.com/you06/chaos-commander/manager/api"
	"github.com/you06/chaos-commander/manager/core"
	"github.com/you06/chaos-commander/manager/scheduler"
	"github.com/you06/chaos-commander/util"
)

var (
	cfg          *config.Config
	printVersion bool
	configPath   string
	logLevel     string
	testMode     bool
	taskID       int
)

func init() {
	flag.StringVar(&configPath, "config-file", "", "path to manager config")
	flag.BoolVar(&printVersion, "V", false, "print version")
	flag.StringVar(&logLevel, "log-level", "info", "log level, support info / warning / debug / error / fatal")
	flag.BoolVar(&testMode, "test", false, "test without scheduler")
	flag.IntVar(&taskID, "task-id", 1, "test with specified ID")
}

func loadConfig() {
	cfg = config.GetGlobalConfig()
	if configPath != "" {
		err := cfg.Load(configPath)
		if err != nil {
			log.Fatalf(errors.ErrorStack(err))
		}
	}
}

func validateConfig() {
	err := util.ValidateListenAddr(cfg.Host)
	if err != nil {
		log.Errorf("verify Listen addr failed %v", err)
		os.Exit(-1)
	}
}

func initLog() {
	log.SetLevelByString(logLevel)
}

func main() {
	flag.Parse()
	if printVersion {
		util.PrintInfo()
		os.Exit(0)
	}

	initLog()
	loadConfig()
	validateConfig()

	if err := cfg.Init(); err != nil {
		log.Error(err)
	}

	mgr, err := core.New(cfg)
	if err != nil {
		log.Fatalf("can't run manager: %v", errors.ErrorStack(err))
	}

	go func() {
		log.Infof("begin to listen %s:%d 😄", cfg.Host, cfg.Port)
		app := iris.New()
		app.Logger().SetLevel(cfg.LogLevel)
		api.CreateRouter(app, mgr)
		app.Run(iris.Addr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
	}()

	sch := scheduler.New(mgr)
	if testMode {
		err = sch.Test(taskID)
		if err != nil {
			log.Errorf("test err %v", err)
		}
	} else {
		sch.Start()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sc
	log.Infof("Got signal %d to exit.", sig)
}
