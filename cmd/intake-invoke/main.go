package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("cmd", "intake-invoke").Logger()
	ctx := logger.WithContext(context.Background())

	flags := parseFlags()
	intakeInvoke, err := quetaro.NewIntakeInvoke(flags.IntakeInvokeOpts)

	if err != nil {
		logger.Fatal().Err(err).Object("flags", flags).
			Msg("failed to create IntakeInvoke struct")
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = intakeInvoke.Start(ctx)

	if err != nil {
		var ewj *quetaro.ErrWithJob

		if errors.As(err, &ewj) {
			logger = logger.With().Str("id", ewj.Id).Str("function_name", ewj.Name).Logger()
		}

		logger.Fatal().Stack().Err(err).Msg("failed to execute intake-invoke")
	}
}
