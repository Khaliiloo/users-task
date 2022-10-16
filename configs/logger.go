package configs

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Logger

func NewLogger(file string) *logrus.Logger {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + file)
		panic(err)
	}

	return &logrus.Logger{
		Out: io.MultiWriter(f, os.Stdout),
		Formatter: &logrus.JSONFormatter{
			PrettyPrint: true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: true,
	}
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}
