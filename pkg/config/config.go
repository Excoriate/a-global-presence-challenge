package config

import (
	"errors"
	"github.com/Excoriate/a-global-presence-challenge/pkg"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Builder struct {
	logger       *zap.Logger
	Env          string
	EnvVars      map[string]string
	Errors       []error
	ChallengeDB  string
	ChallengeDoc string
	ProjectId    string
}

type Config struct {
	Env          string
	Logger       *zap.Logger
	EnvVars      map[string]string
	ChallengeDB  string
	ChallengeDoc string
	ProjectId    string
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
		b.Env = "sandbox"
	}
	return b
}

func (b *Builder) WithScannedEnvVars() *Builder {
	b.EnvVars = make(map[string]string)
	for _, envVar := range os.Environ() {
		b.EnvVars[envVar] = os.Getenv(envVar)
	}

	return b
}

func (b *Builder) WithRequiredConfig() *Builder {
	b.ChallengeDB = os.Getenv("CHALLENGE_DB_NAME")
	b.ChallengeDoc = os.Getenv("CHALLENGE_DOC_NAME")
	b.ProjectId = os.Getenv("PROJECT_ID")

	if b.ChallengeDB == "" {
		b.logger.Error("No CHALLENGE_DB_NAME specified")
		b.Errors = append(b.Errors, errors.New("no CHALLENGE_DB_NAME specified"))
	}

	if b.ChallengeDoc == "" {
		b.logger.Error("No CHALLENGE_DOC_NAME specified")
		b.Errors = append(b.Errors, errors.New("no CHALLENGE_DOC_NAME specified"))
	}

	if b.ProjectId == "" {
		b.logger.Error("No PROJECT_ID specified")
		b.Errors = append(b.Errors, errors.New("no PROJECT_ID specified"))
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
		Env:          b.Env,
		Logger:       b.logger,
		EnvVars:      b.EnvVars,
		ProjectId:    b.ProjectId,
		ChallengeDB:  b.ChallengeDB,
		ChallengeDoc: b.ChallengeDoc,
	}, nil
}
