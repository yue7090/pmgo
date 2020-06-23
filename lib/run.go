package lib

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	runChannel chan string
	stopChannel chan bool
	mainLog    logFunc
	watcherLog logFunc
	runnerLog  logFunc
	buildLog   logFunc
	appLog     logFunc
)

func init() {
	runChannel = make(chan string, 10000)
	stopChannel = make(chan bool)
}

func initLogFuncs() {
	mainLog = newLogFunc("main")
	watcherLog = newLogFunc("watcher")
	runnerLog = newLogFunc("runner")
	buildLog = newLogFunc("build")
	appLog = newLogFunc("app")
} 

func flushEvents() {
	for{
		select {
		case eventName := <-runChannel:
			mainLog("receiving event %s", eventName)
		default:
			return
		}
	}
}

func setEnvVars() {
	os.Setenv("DEV_RUNNER", "1")
	wd, err := os.Getwd()
	if err == nil {
		os.Setenv("RUNNER_WD", wd)
	}
	for k, v := range settings {
		key := strings.ToUpper(fmt.Sprintf("%s%s", envSettingsPrefix, k))
		os.Setenv(key, v)
	}
}

func start() {
	loopIndex := 0
	buildDelay := buildDelay()

	started := false
	go func() {
		for {
			loopIndex++
			mainLog("Waiting (loop %d)...", loopIndex)
			eventName := <-runChannel

			mainLog("receiving first event %s", eventName)
			mainLog("sleeping for %d milliseconds", buildDelay)
			time.Sleep(buildDelay * time.Millisecond)
			mainLog("flushing events")

			flushEvents()

			mainLog("Started! (%d Goroutines)", runtime.NumGoroutine())
			err := removeBuildErrorsLog()
			if err != nil {
				mainLog(err.Error())
			}

			buildFailed := false
			if shouldRebuild(eventName) {
				errorsMessage, ok := build()
				if !ok {
					buildFailed = true
					mainLog("Build Failed:\n %s", errorsMessage)
					if !started {
						os.Exit(1)
					}
					createBuildErrorsLog(errorsMessage)
				}
			}

			if !buildFailed {
				if started {
					stopChannel <- true
				}
				run()
			}


			started = true
			mainLog(strings.Repeat("-", 20))
		}
	}()
}

func Start() {
	initLimit()
	initSettings()
	initLogFuncs()
	initFolders()
	setEnvVars()
	watch()
	start()
	runChannel <- "/"
	<-make(chan int)
}