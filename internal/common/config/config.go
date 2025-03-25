package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      App      `mapstructure:"app"`
		HTTP     HTTP     `mapstructure:"http"`
		Database Database `mapstructure:"database"`
		Logger   Logger   `mapstructure:"logger"`
	}

	App struct {
		Name string `env-required:"true" mapstructure:"name"    env:"APP_NAME"`
	}

	HTTP struct {
		Address      string `env-required:"true" mapstructure:"address"    env:"HTTP_ADDRESS"`
		ReadTimeout  string `env-required:"true" mapstructure:"readTimeout"    env:"HTTP_READ_TIMEOUT"`
		WriteTimeout string `env-required:"true" mapstructure:"writeTimeout"    env:"HTTP_WRITE_TIMEOUT"`
	}

	Database struct {
		Schema   string `env-required:"true" mapstructure:"schema"    env:"DATABASE_SCHEMA"`
		User     string `env-required:"true" mapstructure:"user"    env:"DATABASE_USER"`
		Password string `env-required:"true" mapstructure:"password"    env:"DATABASE_PASSWORD"`
		Port     string `env-required:"true" mapstructure:"port"    env:"DATABASE_PORT"`
		Host     string `env-required:"true" mapstructure:"host"    env:"DATABASE_HOST"`
	}

	Logger struct {
		Handler string `env-required:"true" mapstructure:"handler"    env:"LOGGER_HANDLER"`
		Level   string `env-required:"true" mapstructure:"level"    env:"LOGGER_LEVEL"`
	}
)

// LoadConfig function  
func LoadConfig(configName string) (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName(configName)

	viper.AddConfigPath("./internal/common/config")

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config file because: %w", err)
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to unrmashal config file %s", err)
	}
	return cfg, nil
}
