package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/riabininkf/goragames-assignment/cmd"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		cancelFunc()
	}()

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
