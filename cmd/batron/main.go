package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/takaishi/batron"
)

var Version = "dev"
var Revision = "HEAD"

func init() {
	batron.Version = Version
	batron.Revision = Revision
}

func main() {
	ctx := context.TODO()
	ctx, stop := signal.NotifyContext(ctx, []os.Signal{os.Interrupt}...)
	defer stop()
	if err := batron.RunCLI(ctx, os.Args[1:]); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
