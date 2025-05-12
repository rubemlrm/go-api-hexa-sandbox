package config

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./stubs/config")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestFailLoadConfig(t *testing.T) {
	t.Run("Test failed to load config file", func(t *testing.T) {
		_, err := LoadConfig("invalid_configs")
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to load config file because")
	})

	t.Run("Test failed to unmarshal file", func(t *testing.T) {
		_, err := LoadConfig("invalid_config")
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to unrmashal config file")
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("Test config load without environment variables", func(t *testing.T) {
		cfg, err := LoadConfig("valid_config")
		assert.NoError(t, err)
		assert.Equal(t, "api-gin-template", cfg.App.Name)
		assert.Equal(t, ":8080", cfg.HTTP.Address)
	})

	t.Run("Test config replace with environment variables", func(t *testing.T) {
		err := os.Setenv("APP_NAME", "api-gin-template")
		assert.NoError(t, err)
		err = os.Setenv("HTTP_ADDRESS", ":1234")
		assert.NoError(t, err)
		cfg, err := LoadConfig("valid_config")

		assert.NoError(t, err)
		assert.Equal(t, "api-gin-template", cfg.App.Name)
		assert.Equal(t, ":1234", cfg.HTTP.Address)
	})
}
