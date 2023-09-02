package api_test

import (
	"net/http"
	"testing"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/internal/api"
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

	t.Run("testing right port port", func(t *testing.T) {
		httpConfigs := config.HTTP{
			Address:      "9090",
			ReadTimeout:  "120",
			WriteTimeout: "120",
		}
		err := api.Start(mockHandler{}, httpConfigs)
		assert.Nil(t, err)
	})
}
