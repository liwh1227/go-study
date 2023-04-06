package common

import "github.com/liwh1227/go-common/logger"

var glog *logger.Logger

func init() {
	logger.SetLogConfig(logger.DefaultLogConfig())
	glog = new(logger.Logger)
}

var Log = log()

func log() *logger.Logger {
	return glog
}
