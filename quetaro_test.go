package quetaro_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		os.Setenv("AWS_ACCESS_KEY_ID", "mock_access_key")
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", "mock_secret_key")
	}

	m.Run()
}
