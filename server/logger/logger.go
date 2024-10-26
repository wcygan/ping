package logger

import (
    "go.uber.org/zap"
)

var log *zap.Logger

func Init() {
    var err error
    log, err = zap.NewProduction()
    if err != nil {
        panic(err)
    }
}

func Get() *zap.Logger {
    return log
}
