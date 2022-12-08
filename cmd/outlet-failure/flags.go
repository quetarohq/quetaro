package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/quetarohq/quetaro"
	"github.com/quetarohq/quetaro/cliutil"
)

var (
	version string
)

type Flags struct {
	*quetaro.OutletFailureOpts
}

func init() {
	cmdLine := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	cmdLine.Usage = func() {
		fmt.Fprintf(cmdLine.Output(), "Usage: %s [OPTION]\n", cmdLine.Name())
		cmdLine.PrintDefaults()
	}

	flag.CommandLine = cmdLine
}

func parseFlags() *Flags {
	flags := &Flags{
		OutletFailureOpts: &quetaro.OutletFailureOpts{},
	}

	var dsn string
	flag.StringVar(&flags.OutletFailureOpts.QueueName, "queue", os.Getenv("QTR_OUTLET_FAILURE_QUEUE"), "outlet-failure queue name. use $QTR_OUTLET_FAILURE_QUEUE env")
	flag.StringVar(&dsn, "dsn", os.Getenv("QTR_DATABASE_DSN"), "database DSN. use $QTR_DATABASE_DSN env (e.g. 'postgres://username:password@localhost:5432')")
	flag.IntVar(&flags.NAgents, "nagents", cliutil.GetIntEnv("QTR_OUTLET_FAILURE_NAGENTS", 1), "number of agents. use $QTR_OUTLET_FAILURE_NAGENTS env")
	flag.DurationVar(&flags.Interval, "interval", cliutil.GetDurEnv("QTR_OUTLET_FAILURE_INTERVAL", 1*time.Second), "poll interval. use $QTR_OUTLET_FAILURE_INTERVAL env")
	flag.DurationVar(&flags.ErrInterval, "err-interval", cliutil.GetDurEnv("QTR_OUTLET_FAILURE_ERR_INTERVAL", 1*time.Minute), "error wait interval. use $QTR_OUTLET_FAILURE_ERR_INTERVAL env")
	flag.IntVar(&flags.MaxRecvNum, "max-recv", cliutil.GetIntEnv("QTR_OUTLET_FAILURE_MAX_RECV", 1), "maximum number of received messages. use $QTR_OUTLET_FAILURE_MAX_RECV env")
	flag.StringVar(&flags.AWSRegion, "aws-region", os.Getenv("AWS_REGION"), "AWS region. use $AWS_REGION env")
	flag.StringVar(&flags.AWSEndpointUrl, "aws-endpoint-url", os.Getenv("AWS_ENDPOINT_URL"), "AWS endpoint URL. use $AWS_ENDPOINT_URL env")
	printVer := flag.Bool("version", false, "print version")
	flag.Parse()

	if *printVer {
		cliutil.PrintVersionAndExit(version)
	}

	if flags.OutletFailureOpts.QueueName == "" {
		cliutil.PrintErrorAndExit("'-queue' option is required")
	}

	if dsn == "" {
		cliutil.PrintErrorAndExit("'-dsn' option is required")
	}

	connCfg, err := pgx.ParseConfig(dsn)

	if err != nil {
		cliutil.PrintErrorAndExit("failed to parse database url: %s", err)
	}

	if connCfg.Password == "" {
		connCfg.Password = os.Getenv("PGPASSWORD")
	}

	flags.ConnConfig = connCfg

	return flags
}
