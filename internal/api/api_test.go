package api_test

import (
	"net/http"
	"testing"

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
			ReadTimeout:  "10",
			WriteTimeout: "10",
		}
		err := api.Start(mockHandler{}, httpConfigs)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid port")
	})

	t.Run("testing invalid Read Timeout", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "8080",
			ReadTimeout:  "zxc",
			WriteTimeout: "zxv",
		}
		err := api.Start(mockHandler{}, httpConfigs)
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
		err := api.Start(mockHandler{}, httpConfigs)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Error validating configuration: WriteTimeout")
		assert.IsType(t, err, &gin.HttpConfigurationError{})
	})
}
