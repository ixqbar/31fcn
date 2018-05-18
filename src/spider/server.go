package spider

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/robfig/cron"
)

var SpiderGlobalVars = NewTGlobalVars()

func Run() {
	SpiderGlobalVars.Init();

	f31proxy := NewF31Proxy()
	if GTSpiderConfig.PullOnStartup {
		f31proxy.Pull()
	}

	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)

	c := cron.New()
	c.AddFunc(GTSpiderConfig.Schedule, func() {
		index := SpiderGlobalVars.Index()
		Logger.Printf("schedule[%d] start execute", index)
		defer Logger.Printf("schedule[%d] stop execute", index)

		f31proxy.Pull()
	})
	c.Start()

	select {
		case <-stopSignal:
			Logger.Print("catch exit signal")
	}

	c.Stop()
	SpiderGlobalVars.Clear()

	Logger.Print("proxy spider exit")
}