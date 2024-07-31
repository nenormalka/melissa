package config

import (
	"fmt"

	types "github.com/nenormalka/melissa/appinfo"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		ReleaseID string
		Env       string `envconfig:"ENV" default:"development" required:"true" yaml:"env"`
		LogLevel  string `envconfig:"LOG_LEVEL" default:"info" yaml:"log_level"`
		AppName   string `envconfig:"APP_NAME" yaml:"app_name"`
	}
)

func NewConfig(info *types.AppInfo) (*Config, error) {
	types.SetAppInfo(info)

	cfg := &Config{}
	cfg.ReleaseID = types.GetAppVersion()

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("parse env config err %w", err)
	}

	return cfg, nil
}
