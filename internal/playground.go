package main

import (
	"io/ioutil"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/syllabix/logger/mode"

	"github.com/syllabix/logger"
)

/****
	NOTE:
	This file is simple a logger "playground" - experiments and/or examples worth perserving should be added
	here
****/

func main() {

	logger.Configure(
		logger.AppName("test-logger"),
		logger.Mode(mode.Production),
		logger.JSONWriter(os.Stderr),
	)

	log := logger.New().Sugar()

	then := time.Now()
	time.Sleep(100 * time.Millisecond)
	log.Errorw("this is a time",
		"now", time.Now(),
		"since", time.Since(then))

	// update config for all loggers
	logger.Configure(
		logger.AppName("foo-bazz"),
		logger.Mode(mode.Development),
		logger.JSONWriter(ioutil.Discard),
	)

	// perf logger
	plog := logger.New()

	time.Sleep(100 * time.Millisecond)
	plog.Info("this is a time",
		zap.Time("now", time.Now()),
		zap.Duration("since", time.Since(then)))
	plog.Warn("this is a time",
		zap.Time("now", time.Now()),
		zap.Duration("since", time.Since(then)))
	plog.Error("this is a time",
		zap.Time("now", time.Now()),
		zap.Duration("since", time.Since(then)))
}
