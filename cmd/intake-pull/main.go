package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/winebarrel/qtr"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("cmd", "intake-pull").Logger()
	ctx := logger.WithContext(context.Background())

	flags := parseFlags()
	intakePull, err := qtr.NewIntakePull(flags.IntakePullOpts)

	if err != nil {
		logger.Fatal().Err(err).Object("flags", flags).
			Msg("failed to create IntakePull struct")
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = intakePull.Start(ctx)

	if err != nil {
		var ewj *qtr.ErrWithJob

		if errors.As(err, &ewj) {
			logger = logger.With().Str("id", ewj.Id).Str("function_name", ewj.Name).Logger()
		}

		logger.Fatal().Stack().Err(err).Msg("failed to execute intake-pull")
	}
}
