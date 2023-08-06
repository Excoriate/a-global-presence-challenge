package hackattic

import (
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
)

type AttemptBuilder struct {
	Cfg       *config.Config
	Challenge *Challenge
	Errors    []error
	Token     string
	Result    *Result
}

type Attempt struct {
	Cfg       *config.Config
	Challenge *Challenge
	Token     string
	Result    *Result
}

type Result struct {
	CountryResponse string
	IsSuccessful    bool
}

func (b *AttemptBuilder) WithChallenge(challenge *Challenge) *AttemptBuilder {
	b.Challenge = challenge
	return b
}

func (b *AttemptBuilder) WithToken(token string) *AttemptBuilder {
	// Grab the token from the challenge.
	return nil
}

func (b *AttemptBuilder) WithExistingToken() *AttemptBuilder {
	if b.Cfg.Env == "dev" {
		b.Token = b.Cfg.EnvVars["PRESENCE_TOKEN"]
		if b.Token == "" {
			b.Cfg.Logger.Error("Token cannot be empty")
			b.Errors = append(b.Errors, fmt.Errorf("token cannot be empty"))
			return b
		}

		b.Cfg.Logger.Info("Using token from environment variables")
	}

	return b
}
