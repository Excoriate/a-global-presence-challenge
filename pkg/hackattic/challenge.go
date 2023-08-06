package hackattic

import (
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Challenge struct {
	Headers             []string
	Endpoint            string // Normally, the challenge name?
	Token               string
	APIGetPresenceToken string // We get the token.
	APICountryCheck     string // We check the country.
	APISubmitSolution   string // We submit the solution.
	HttpClient          *http.Client
}

type ChallengeBuilder struct {
	Cfg               *config.Config
	Base              string
	APIBase           string
	APIInit           string // We get the token.
	APICountryCheck   string // We check the country.
	APISubmitSolution string // We submit the solution.
	Headers           []string
	Errors            []error
	Token             string
	HttpClient        *http.Client
}

func New(cfg *config.Config) *ChallengeBuilder {
	return &ChallengeBuilder{
		Cfg:     cfg,
		Base:    "https://hackattic.com",
		APIBase: "challenges/a_global_presence",
	}
}

func (b *ChallengeBuilder) WithCountryHeaders(country string) *ChallengeBuilder {
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

func (b *ChallengeBuilder) WithNewToken() *ChallengeBuilder {
	// TODO: Implement it later. Perhaps in GCP,
	// it'd go and check a FireStore DB and retrieve it from there?
	return nil
}

func (b *ChallengeBuilder) WithHTTPClient() *ChallengeBuilder {
	client := &http.Client{}

	// Prepare headers into http.Header
	headers := make(http.Header)
	for _, h := range b.Headers {
		headerParts := strings.SplitN(h, ":", 2)
		if len(headerParts) != 2 {
			b.Errors = append(b.Errors, fmt.Errorf("invalid header format: %s", h))
			return b
		}
		headers.Set(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
	}

	// Replacing the client's Transport function to inject headers into every request
	client.Transport = &headerTransport{
		Transport: http.DefaultTransport,
		Headers:   headers,
	}

	// storing the configured http client in the builder
	b.HttpClient = client
	return b
}

type headerTransport struct {
	Transport http.RoundTripper
	Headers   http.Header
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// copy headers to the new request to avoid modifying original request
	for k, vv := range t.Headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	return t.Transport.RoundTrip(req)
}

func (b *ChallengeBuilder) Build() (*Challenge, error) {
	if len(b.Errors) > 0 {
		for _, err := range b.Errors {
			b.Cfg.Logger.Error("Error(s) found in builder", zap.Error(err))
		}
		return nil, fmt.Errorf("error(s) found in builder")
	}

	url := &Challenge{
		Headers:             b.Headers,
		Endpoint:            b.APIBase,
		Token:               b.Token,
		APIGetPresenceToken: fmt.Sprintf("%s/%s/problem?access_token=%s", b.Base, b.APIBase, b.Token),
		APICountryCheck:     fmt.Sprintf("%s/_/presence/%s", b.Base, b.Token),
		APISubmitSolution:   fmt.Sprintf("%s/%s/solve?access_token=%s", b.Base, b.APIBase, b.Token),
		HttpClient:          b.HttpClient,
	}

	return url, nil
}
