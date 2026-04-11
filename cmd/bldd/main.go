package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/serrasma/bldd/internal/application/bldd"
	"github.com/serrasma/bldd/internal/config"
	"github.com/serrasma/bldd/internal/domain"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	libraries := flag.String("lib", "", "libraries to find")
	directories := flag.String("dir", "", "directories with binaries")
	reportPath := flag.String("report", "bldd_report.md", "report file path")
	workers := flag.Int("worker", 4, "number of workers")

	flag.Parse()

	cfg := config.New(*libraries, *directories, *reportPath, *workers)
	walker := domain.NewWalker(cfg)
	reporter := domain.NewReporter(cfg)

	errs := make(chan error)
	bldd := bldd.New(walker, reporter, errs)
	bldd.Run(ctx)

	select {
	case <-ctx.Done():
		stop()
	case err := <-errs:
		if err != nil {
			slog.Error("bldd run error", slog.String("err", err.Error()))
		}
	}
}
