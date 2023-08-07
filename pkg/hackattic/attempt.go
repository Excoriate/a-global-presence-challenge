package hackattic

import (
	"encoding/json"
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
)

type AttemptBuilder struct {
	Cfg           *config.Config
	Challenge     *Challenge
	Errors        []error
	PresenceToken string
	Result        *Result
}

type Attempt struct {
	Cfg           *config.Config
	Challenge     *Challenge
	PresenceToken string
	Result        *Result
}

type Result struct {
	CountryResponse string
	IsSuccessful    bool
}

func (b *AttemptBuilder) WithChallenge(challenge *Challenge) *AttemptBuilder {
	b.Challenge = challenge
	return b
}

type presenceTokenResponse struct {
	PresenceToken string `json:"presence_token"`
}

func (b *AttemptBuilder) WithNewPresenceToken() *AttemptBuilder {
	// Grab the token from the challenge.
	accessAPI := b.Challenge.APIGetPresenceToken
	b.Cfg.Logger.Info(fmt.Sprintf("Attempting to get presence token from %s", accessAPI))

	resp, err := b.Challenge.HttpClient.Get(accessAPI)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error getting presence token: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	body := resp.Body
	defer body.Close()

	// Decode the response.
	var tokenResponse presenceTokenResponse
	err = json.NewDecoder(body).Decode(&tokenResponse)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error decoding presence token response: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	b.PresenceToken = tokenResponse.PresenceToken
	b.Challenge.APICountryCheck = fmt.Sprintf("%s/%s", b.Challenge.APICountryCheck,
		b.PresenceToken)
	b.Cfg.EnvVars["PRESENCE_TOKEN"] = tokenResponse.PresenceToken
	b.Cfg.Logger.Info(fmt.Sprintf("Presence token set to %s", b.PresenceToken))

	return b
}

func (b *AttemptBuilder) Build() (*Attempt, error) {
	if len(b.Errors) > 0 {
		b.Cfg.Logger.Error(fmt.Sprintf("Cannot create attempt with errors: %v", b.Errors))
		for _, err := range b.Errors {
			b.Cfg.Logger.Error(fmt.Sprintf("Error: %s", err.Error()))
		}

		return nil, fmt.Errorf("cannot create attempt with errors")
	}

	return &Attempt{
		Cfg:           b.Cfg,
		Challenge:     b.Challenge,
		PresenceToken: b.PresenceToken,
		Result:        &Result{},
	}, nil
}

func NewAttempt(cfg *config.Config) *AttemptBuilder {
	return &AttemptBuilder{
		Cfg:    cfg,
		Errors: []error{},
	}
}
