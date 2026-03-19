package main

import (
	"context"
	"fmt"
	"os"

	"github.com/maxanderson95/md2pdf/internal/app"
	"github.com/maxanderson95/md2pdf/internal/version"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, version.Info()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if exitCoder, ok := err.(interface{ ExitCode() int }); ok {
			os.Exit(exitCoder.ExitCode())
		}
		os.Exit(1)
	}
}
