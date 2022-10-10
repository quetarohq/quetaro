package cliutil

import (
	"flag"
	"fmt"
	"os"
)

func PrintVersionAndExit(version string) {
	v := version

	if v == "" {
		v = "<nil>"
	}

	fmt.Fprintln(flag.CommandLine.Output(), v)
	os.Exit(0)
}

func PrintErrorAndExit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
