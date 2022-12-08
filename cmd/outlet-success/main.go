package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/quetarohq/quetaro"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("cmd", "outlet-success").Logger()
	ctx := logger.WithContext(context.Background())
	flags := parseFlags()
	outletSuccess, err := quetaro.NewOutletSuccess(flags.OutletSuccessOpts)

	if err != nil {
		logger.Fatal().Err(err).Object("flags", flags).
			Msg("failed to create OutletSuccess struct")
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = outletSuccess.Start(ctx)

	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("failed to execute outlet-success")
	}
}
