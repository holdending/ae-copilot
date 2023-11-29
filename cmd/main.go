package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/LiveRamp/ae-copilot/config"
	"github.com/LiveRamp/ae-copilot/pkg/libs/logger"
	"github.com/LiveRamp/ae-copilot/routers"
	"github.com/LiveRamp/ae-copilot/scan"
	"github.com/astaxie/beego/logs"
)

func main() {
	runtime.GOMAXPROCS(128)
	scan.AsyncRunning()
	logger.Initialize(config.Agent.LogType, config.Agent.LogConf, config.Agent.LogLevel, config.Agent.SendgridConf)

	go func() {
		logs.Error(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var HTTPAddr string
	flag.StringVar(&HTTPAddr, "addr", "", "")
	flag.Parse()

	if HTTPAddr == "" {
		HTTPAddr = ":" + config.Agent.HTTPPort
	}
	router := routers.NewRouter()
	logs.Error(http.ListenAndServe(HTTPAddr, router))
}
