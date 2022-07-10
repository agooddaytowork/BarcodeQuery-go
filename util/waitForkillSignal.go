package util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForKillSignal() {
	s := make(chan os.Signal, 1)
	signal.Notify(s,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-s
}
