package cliutil

import (
	"os"
	"strconv"
	"time"
)

func GetIntEnv(name string, defval int) int {
	s := os.Getenv(name)

	if s == "" {
		return defval
	}

	i, err := strconv.Atoi(s)

	if err != nil {
		PrintErrorAndExit("failed to parse $%s: %s", name, err)
	}

	return i
}

func GetDurEnv(name string, defval time.Duration) time.Duration {
	s := os.Getenv(name)

	if s == "" {
		return defval
	}

	d, err := time.ParseDuration(s)

	if err != nil {
		PrintErrorAndExit("failed to parse $%s: %s", name, err)
	}

	return d
}
