package main

import (
	"github.com/sethpollack/bogie/cmd"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/logging"
)

func init() {
	// change sops log level to error
	for _, log := range logging.Loggers {
		log.SetLevel(logrus.ErrorLevel)
	}
}

func main() {
	cmd.Execute()
}
