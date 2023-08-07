package hackattic

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/Excoriate/a-global-presence-challenge/pkg/utils"
	"io"
)

type Shooter struct {
	Id         string
	AttemptId  string
	Status     string
	Cfg        *config.Config
	Challenge  *Challenge
	Attempt    *Attempt
	Errors     []error
	AttemptDoc *Document
	// The response back from hackattic
	Response *ShooterResponse
}

type ShooterResponse struct {
	HackatticCountryCheck string
	IsSuccessful          bool
	Status                string
	Error                 error
}

type ShooterBuilder struct {
	Cfg        *config.Config
	Challenge  *Challenge
	Attempt    *Attempt
	Id         string
	Errors     []error
	AttemptDoc *Document
	// The response back from hackattic
	Response *ShooterResponse
}

func NewShooter() *ShooterBuilder {
	return &ShooterBuilder{}
}

func (b *ShooterBuilder) WithConfig(cfg *config.Config) *ShooterBuilder {
	b.Cfg = cfg
	return b
}

func (b *ShooterBuilder) WithChallenge(challenge *Challenge) *ShooterBuilder {
	b.Challenge = challenge
	return b
}

func (b *ShooterBuilder) WithAttempt(attempt *Attempt) *ShooterBuilder {
	b.Attempt = attempt
	return b
}

type AttemptRegisterDoc struct {
	AttemptID             string `firestore:"attempt_id"`
	PresenceToken         string `firestore:"presenceToken"`
	IsCompleted           bool   `firestore:"isCompleted"`
	Status                string `firestore:"status"`
	HackatticCountryCheck string `firestore:"hackatticCountryCheck"`
}

type Document struct {
	AttemptID        string               `firestore:"attempt_id"`
	AccessToken      string               `firestore:"accessToken"`
	PresenceToken    string               `firestore:"presenceToken"`
	Attempts         []AttemptRegisterDoc `firestore:"attempts"`
	Status           string               `firestore:"status"`
	NumberOfAttempts int                  `firestore:"numberOfAttempts"`
	Result           string               `firestore:"result"`
}

func (b *ShooterBuilder) WithPassedPresenceToken(token string) *ShooterBuilder {
	if token == "" {
		b.Cfg.Logger.Error("Empty presence token passed")
		b.Errors = append(b.Errors, fmt.Errorf("empty presence token passed"))
		return b
	}

	b.Attempt.PresenceToken = token
	return b
}

func (b *ShooterBuilder) WithPresenceTokenFromDb() *ShooterBuilder {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, b.Cfg.ProjectId)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error creating firestore client: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	defer client.Close()

	docs := client.Collection(b.Cfg.ChallengeDoc).Where("status", "!=", "completed").OrderBy("status", firestore.Desc).Documents(ctx)
	doc, err := docs.Next()
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error getting document from firestore: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	var attempt Document
	err = doc.DataTo(&attempt)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error converting firestore document to struct: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	b.Attempt.PresenceToken = attempt.PresenceToken
	b.Cfg.Logger.Info(fmt.Sprintf("Presence token from firestore: %s", b.Attempt.PresenceToken))
	b.AttemptDoc = &attempt

	return b
}

func (b *ShooterBuilder) WithCountryCheck() *ShooterBuilder {
	presenceToken := b.Attempt.PresenceToken
	resp, err := b.Challenge.HttpClient.Get(fmt.Sprintf("%s/%s", b.Challenge.APICountryCheck, presenceToken))

	if presenceToken == "" {
		b.Cfg.Logger.Error(fmt.Sprintf("Error getting presence token from firestore: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error getting country check from hackattic: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error reading response body: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	bodyString := string(body)
	b.Cfg.Logger.Info(fmt.Sprintf("Country check response: %s", bodyString))

	b.Response = &ShooterResponse{
		HackatticCountryCheck: bodyString,
		IsSuccessful:          true,
		Status:                "checked",
		Error:                 nil,
	}

	return b
}

func (b *ShooterBuilder) Complete() (*Shooter, []error) {
	if len(b.Errors) > 0 {
		return nil, b.Errors
	}

	return &Shooter{
		Id:         utils.GetUUID(),
		AttemptId:  utils.GetUUID(),
		Status:     "pending",
		Cfg:        b.Cfg,
		Challenge:  b.Challenge,
		Attempt:    b.Attempt,
		AttemptDoc: b.AttemptDoc,
		Response:   b.Response,
	}, nil
}

type ShooterExecutioner interface {
	Register() *ShooterResponse
}

func (s *Shooter) Register() *ShooterResponse {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, s.Cfg.ProjectId)
	if err != nil {
		s.Cfg.Logger.Error(fmt.Sprintf("Error creating firestore client: %s", err.Error()))
		s.Errors = append(s.Errors, err)
		return nil
	}

	// Update the attempt in the firestore
	s.AttemptDoc.Attempts = append(s.AttemptDoc.Attempts, AttemptRegisterDoc{
		AttemptID:             utils.GetUUID(),
		HackatticCountryCheck: s.Response.HackatticCountryCheck,
		IsCompleted:           true,
		PresenceToken:         s.Attempt.PresenceToken,
		Status:                "hackattic-checked-ok",
	})

	defer client.Close()

	_, err = client.Collection(s.Cfg.ChallengeDoc).Doc(s.AttemptDoc.AttemptID).Set(ctx, s.AttemptDoc)

	if err != nil {
		s.Cfg.Logger.Error(fmt.Sprintf("Error updating document in firestore: %s", err.Error()))
		s.Errors = append(s.Errors, err)
		return nil
	}

	return s.Response
}
