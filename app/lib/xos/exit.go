package xos

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}
