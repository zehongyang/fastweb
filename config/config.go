package config

import (
	"flag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testing"
)

var (
	gData   []byte
	debug   bool
	GLogger *zap.Logger
	err     error
)

func init() {
	gData, err = os.ReadFile("./application.yml")
	if err != nil {
		log.Fatal("ReadFile Err:", err)
	}
	flag.BoolVar(&debug, "d", false, "true is debug false is production")
	testing.Init()
	flag.Parse()
	initLogger()
}

func Load(obj any) error {
	return yaml.Unmarshal(gData, obj)
}

func IsDebug() bool {
	return debug
}

func initLogger() {
	var cfg zap.Config
	if !IsDebug() {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.DisableStacktrace = true
	GLogger, err = cfg.Build()
	if err != nil {
		log.Fatal("initLogger Err:", err)
	}
	GLogger = GLogger.WithOptions(zap.AddCallerSkip(1))
}
