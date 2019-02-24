package errortracking

import (
	raven "github.com/getsentry/raven-go"
)

type Config struct {
	Version     string `envconfig:"VERSION"`
	Environment string `envconfig:"ENVIRONMENT"`
	RavenDSN    string `envconfig:"RAVEN_DSN"`
}

func Init(cfg *Config) error {
	if cfg.RavenDSN == "" {
		return nil
	}

	err := raven.SetDSN(cfg.RavenDSN)
	if err != nil {
		return err
	}
	raven.SetEnvironment(cfg.Environment)
	raven.SetRelease(cfg.Version)
	return nil
}
