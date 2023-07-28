package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App  App  `mapstructure:"app"`
		HTTP HTTP `mapstructure:"http"`
	}

	App struct {
		Name string `env-required:"true" mapstructure:"name"    env:"APP_NAME"`
	}

	HTTP struct {
		Address string `env-required:"true" mapstructure:"address"    env:"HTTP_ADDRESS"`
	}
)

// LoadConfig function  
func LoadConfig(configName string) (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName(configName)

	viper.AddConfigPath("./config")

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config file because: %w", err)
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to unrmashal config file %s", err)
	}
	return cfg, nil
}
