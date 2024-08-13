package service

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/arlettebrook/serv00-ct8/configs"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	switch strings.ToLower(configs.Cfg.LogLevel) {
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
	Logger.SetFormatter(&nested.Formatter{
		TimestampFormat: time.DateTime,
		ShowFullLevel:   true,
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			return fmt.Sprintf(" %s:%d", filepath.Base(f.File), f.Line)
		},
	})
	Logger.SetReportCaller(true)
}
