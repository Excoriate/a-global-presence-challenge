package hackattic

import (
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

type Challenge struct {
	Headers  []string
	Endpoint string // Normally, the challenge name?
	// Tokens
	AccessToken string // Allow to retrieve a presence token
	// APIs
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
	AccessToken       string // Allow to retrieve a presence token
	PresenceToken     string // Allow to perform challenge attempts.
	HttpClient        *http.Client
}

func New(cfg *config.Config) *ChallengeBuilder {
	return &ChallengeBuilder{
		Cfg:     cfg,
		Base:    "https://hackattic.com",
		APIBase: "challenges/a_global_presence",
	}
}

func (b *ChallengeBuilder) WithAccessToken() *ChallengeBuilder {
	token := os.Getenv("ACCESS_TOKEN")

	if token == "" {
		b.Cfg.Logger.Error("Access token cannot be empty")
		b.Errors = append(b.Errors, fmt.Errorf("access token cannot be empty"))
		return b
	}

	b.Cfg.EnvVars["ACCESS_TOKEN"] = token
	b.AccessToken = token
	b.Cfg.Logger.Info(fmt.Sprintf("Access token set to %s", token))
	return b
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

func (b *ChallengeBuilder) WithHeaders(headers []string) *ChallengeBuilder {
	headers = append(headers, "Content-Type: application/json")

	if len(headers) > 0 {
		headers = append(headers, b.Headers...)
	}

	b.Headers = headers
	return b
}

func (b *ChallengeBuilder) Build() (*Challenge, []error) {
	if len(b.Errors) > 0 {
		for _, err := range b.Errors {
			b.Cfg.Logger.Error("Error(s) found in builder", zap.Error(err))
		}

		return nil, b.Errors
	}

	url := &Challenge{
		Headers:  b.Headers,
		Endpoint: b.APIBase,
		APIGetPresenceToken: fmt.Sprintf("%s/%s/problem?access_token=%s", b.Base, b.APIBase,
			b.AccessToken),
		APICountryCheck:   fmt.Sprintf("%s/_/presence", b.Base),
		APISubmitSolution: fmt.Sprintf("%s/%s/solve?access_token=%s", b.Base, b.APIBase, b.AccessToken),
		HttpClient:        b.HttpClient,
		AccessToken:       b.AccessToken,
	}

	return url, nil
}
