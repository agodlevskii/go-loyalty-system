package logs

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	Logger = logger
}
