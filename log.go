package main

import (
	"os/user"
	ol "log"

	"github.com/antigloss/go/logger"
)

var logLevels = map[string]logger.LogDest{"file": logger.LogDestFile, "console": logger.LogDestConsole, "none": logger.LogDestNone}
var logLevelKey string

var (
	homeDir string
	faciDir string
	logDir string
)

var log *logger.Logger

func setHomeDir() {
	u, err := user.Current()
	if err != nil {
		ol.Fatal(err)
	}
	homeDir = u.HomeDir + "/"
}

func setFaciDir() {
	faciDir = homeDir + ".faci/"
}

func setLogDir() {
	logDir = faciDir + "log/"
}

func newLog() *logger.Logger {
	config := &logger.Config {
		LogDir: logDir,
		LogFileMaxSize: 1,
		LogFileMaxNum: 10,
		LogFileNumToDel: 1,
		LogLevel: logger.LogLevelInfo,
		LogDest: logLevels[logLevelKey],
		Flag: logger.ControlFlagLogDate | logger.ControlFlagLogLineNum,
	}
	log, err := logger.New(config)
	if err != nil {
		ol.Fatal(err)
	}

	return log
}
