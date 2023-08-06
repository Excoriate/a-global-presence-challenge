package main

import (
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/Excoriate/a-global-presence-challenge/pkg/hackattic"
	"io"
)

func main() {
	client := config.New()
	// Get configuration.
	cfg, err := client.
		WithEnv("").
		WithDotEnv().
		Build()

	if err != nil {
		panic(err)
	}

	// Get Challenge
	h, err := hackattic.New(cfg).
		WithHTTPClient().
		Build()

	if err != nil {
		panic(err)
	}

	resp, err := h.HttpClient.Get(h.APICountryCheck)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Get the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body.
	println(string(body))
}
