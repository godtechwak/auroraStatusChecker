package logger_util

import (
	"log"
	"os"
)

func NewLogger(filename string) (*log.Logger, func(), error) {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(logFile, "", log.LstdFlags)
	cleanup := func() {
		logFile.Close()
	}

	return logger, cleanup, nil
}
