package api_test

import (
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
	"github.com/rubemlrm/go-api-bootstrap/pkg/gin"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct{}

func (m mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestStart(t *testing.T) {

	t.Run("testing wrong port", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "99999999999",
			ReadTimeout:  "1",
			WriteTimeout: "1",
		}
		srv, err := api.NewServer(mockHandler{}, httpConfigs)
		assert.Nil(t, err)
		go func(error) {
			time.Sleep(1 * time.Second)
			p, err := os.FindProcess(syscall.Getpid())
			assert.Nil(t, err)
			err = p.Signal(syscall.SIGINT)
			assert.Nil(t, err)
		}(err)
		err = srv.Start()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "listen tcp: address 99999999999: invalid port")
	})

	t.Run("testing invalid Read Timeout", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "8080",
			ReadTimeout:  "zxc",
			WriteTimeout: "zxv",
		}
		server, err := api.NewServer(mockHandler{}, httpConfigs)
		assert.IsType(t, server, api.Server{Server: (*http.Server)(nil)})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Error validating configuration: ReadTimeout")
		assert.IsType(t, err, &gin.HttpConfigurationError{})
	})

	t.Run("testing invalid Write Timeout", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "8080",
			ReadTimeout:  "123",
			WriteTimeout: "zxv",
		}
		server, err := api.NewServer(mockHandler{}, httpConfigs)
		assert.IsType(t, server, api.Server{Server: (*http.Server)(nil)})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Error validating configuration: WriteTimeout")
		assert.IsType(t, err, &gin.HttpConfigurationError{})
	})

	t.Run("testing Server Start", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "8080",
			ReadTimeout:  "1",
			WriteTimeout: "1",
		}
		srv, err := api.NewServer(mockHandler{}, httpConfigs)
		assert.NoError(t, err)

		go func() {
			time.Sleep(1 * time.Second)
			p, err := os.FindProcess(syscall.Getpid())
			assert.Nil(t, err)
			err = p.Signal(syscall.SIGINT)
			assert.Nil(t, err)

		}()
		err = srv.Start()
		assert.Nil(t, err)
	})
}
