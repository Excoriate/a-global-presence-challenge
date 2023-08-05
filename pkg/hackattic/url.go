package hackattic

import (
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"go.uber.org/zap"
	"strings"
)

type URL struct {
	Headers  []string
	Endpoint string
	Token    string
}

type URLBuilder struct {
	Cfg       *config.Config
	Base      string
	API       string
	Headers   []string
	Challenge string
	Errors    []error
	Token     string
}

func NewURLBuilder(cfg *config.Config) *URLBuilder {
	return &URLBuilder{
		Cfg:       cfg,
		Base:      "https://hackattic.com",
		API:       "challenges",
		Challenge: "",
	}
}

func (b *URLBuilder) WithChallenge(name string) *URLBuilder {
	if name == "" {
		b.Cfg.Logger.Error("Challenge cannot be empty")
		b.Errors = append(b.Errors, fmt.Errorf("challenge cannot be empty"))
		return b
	}

	b.Challenge = strings.TrimSpace(name)
	b.Cfg.Logger.Info(fmt.Sprintf("Building URL for challenge %s", b.Challenge))
	return b
}

func (b *URLBuilder) WithCountryHeaders(country string) *URLBuilder {
	if country == "" {
		b.Cfg.Logger.Error("Country cannot be empty")
		b.Errors = append(b.Errors, fmt.Errorf("country cannot be empty"))
		return b
	}

	countryId := strings.ToUpper(country)
	headers := GetHeaders()

	for _, header := range headers {
		if header.Country != countryId {
			continue
		}

		b.Headers = append(b.Headers, header.Headers...)
		b.Cfg.Logger.Info(fmt.Sprintf("Added headers for country %s", countryId))
	}

	return b
}

func (b *URLBuilder) WithToken() *URLBuilder {
	if b.Cfg.EnvVars["PRESENCE_TOKEN"] == "" {
		b.Cfg.Logger.Error("Presence token not found")
		b.Errors = append(b.Errors, fmt.Errorf("presence token not found"))
		return b
	}

	b.Token = b.Cfg.EnvVars["PRESENCE_TOKEN"]
	return b
}

func (b *URLBuilder) Build() (*URL, error) {
	if len(b.Errors) > 0 {
		for _, err := range b.Errors {
			b.Cfg.Logger.Error("Error(s) found in builder", zap.Error(err))
		}
		return nil, fmt.Errorf("error(s) found in builder")
	}

	url := &URL{
		Headers:  b.Headers,
		Endpoint: fmt.Sprintf("%s/%s/%s", b.Base, b.API, b.Challenge),
		Token:    b.Token,
	}

	return url, nil
}
