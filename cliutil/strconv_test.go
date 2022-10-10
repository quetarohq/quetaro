package cliutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro/cliutil"
)

func Test_GetIntEnv(t *testing.T) {
	t.Setenv("__QuetaroTest_GetIntEnv__", "3")
	assert := assert.New(t)
	i := cliutil.GetIntEnv("__QuetaroTest_GetIntEnv__", 1)
	assert.Equal(3, i)
}

func Test_GetIntEnv_Default(t *testing.T) {
	assert := assert.New(t)
	i := cliutil.GetIntEnv("__QuetaroTest_GetIntEnv__", 1)
	assert.Equal(1, i)
}

func Test_GetDurEnv(t *testing.T) {
	t.Setenv("__QuetaroTest_Test_GetDurEnv__", "3s")
	assert := assert.New(t)
	i := cliutil.GetDurEnv("__QuetaroTest_Test_GetDurEnv__", 1*time.Second)
	assert.Equal(3*time.Second, i)
}

func Test_GetDurEnv_Default(t *testing.T) {
	assert := assert.New(t)
	i := cliutil.GetDurEnv("__QuetaroTest_Test_GetDurEnv__", 1*time.Second)
	assert.Equal(1*time.Second, i)
}
