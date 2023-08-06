package config

import (
	"errors"
	"github.com/Excoriate/a-global-presence-challenge/pkg"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Builder struct {
	logger  *zap.Logger
	Env     string
	EnvVars map[string]string
	Errors  []error
}

type Config struct {
	Env     string
	Logger  *zap.Logger
	EnvVars map[string]string
}

func New() *Builder {
	return &Builder{
		logger: pkg.NewLogger(),
	}
}

func (b *Builder) WithEnv(env string) *Builder {
	b.Env = env
	if b.Env == "" {
		b.logger.Info("No environment specified, defaulting to prod")
		b.Env = "prod"
	}
	return b
}

func (b *Builder) WithDotEnv() *Builder {
	if _, err := os.Stat(".env"); err == nil {
		b.logger.Info("Found .env file, loading environment variables")
		err := godotenv.Load()
		if err != nil {
			b.logger.Error("Error loading .env file", zap.Error(err))
			b.Errors = append(b.Errors, err)
		} else {
			b.EnvVars, err = godotenv.Read()
			if err != nil {
				b.logger.Error("Error reading .env file", zap.Error(err))
				b.Errors = append(b.Errors, err)
			}
		}
	} else {
		b.logger.Info("No .env file found. Most likely running in production.")
	}

	return b
}

func (b *Builder) Build() (*Config, error) {
	if len(b.Errors) > 0 {
		for _, err := range b.Errors {
			b.logger.Error("Error(s) found in builder", zap.Error(err))
		}

		return nil, errors.New("error(s) found in builder")
	}

	return &Config{
		Env:     b.Env,
		Logger:  b.logger,
		EnvVars: b.EnvVars,
	}, nil
}
