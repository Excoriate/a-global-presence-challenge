package hackattic

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/Excoriate/a-global-presence-challenge/pkg/utils"
)

type Shooter struct {
	Id        string
	AttemptId string
	Status    string
	Cfg       *config.Config
	Challenge *Challenge
	Attempt   *Attempt
	Errors    []error
}

type ShooterResponse struct {
	HackatticCountryCheck string
	IsSuccessful          bool
	Status                string
	Error                 error
}

type ShooterBuilder struct {
	Cfg       *config.Config
	Challenge *Challenge
	Attempt   *Attempt
	Id        string
	Errors    []error
}

func (b *ShooterBuilder) NewShooter() *Shooter {
	return &Shooter{}
}

type AttemptShooterRegister struct {
	AttemptId             string `firestore:"attempt_id"`
	IsCompleted           bool   `firestore:"is_completed"`
	Status                string `firestore:"status"`
	PresenceToken         string `firestore:"presence_token"`
	HackatticCountryCheck string `firestore:"hackattic_country_check"`
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

	doc, err := client.Collection(b.Cfg.ChallengeDoc).Doc(b.Id).Get(ctx)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error getting document from firestore: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	var attempt AttemptShooterRegister
	err = doc.DataTo(&attempt)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error converting firestore document to struct: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	b.Attempt.PresenceToken = attempt.PresenceToken
	b.Cfg.Logger.Info(fmt.Sprintf("Presence token from firestore: %s", b.Attempt.PresenceToken))

	return b // returning modified builder
}

func (b *ShooterBuilder) WithNewAttempt() *ShooterBuilder {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, b.Cfg.ProjectId)
	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error creating firestore client: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	defer client.Close()

	newAttempt := map[string]interface{}{
		"attempt_id":            utils.GetUUID(),
		"presenceToken":         b.Attempt.PresenceToken,
		"isCompleted":           false,
		"status":                "trigger-init",
		"hackatticCountryCheck": "",
	}

	_, err = client.Collection(b.Cfg.ChallengeDoc).Doc(b.Id).Update(ctx, []firestore.Update{
		{
			Path:  "attempts",
			Value: firestore.ArrayUnion(newAttempt),
		},
	})

	if err != nil {
		b.Cfg.Logger.Error(fmt.Sprintf("Error updating document in firestore: %s", err.Error()))
		b.Errors = append(b.Errors, err)
		return b
	}

	return b
}

func (b *ShooterBuilder) Complete() (*Shooter, []error) {
	if len(b.Errors) > 0 {
		return nil, b.Errors
	}

	return &Shooter{
		Id:        utils.GetUUID(),
		AttemptId: utils.GetUUID(),
		Status:    "pending",
		Cfg:       b.Cfg,
		Challenge: b.Challenge,
		Attempt:   b.Attempt,
	}, nil
}

type ShooterExecutioner interface {
	Execute() *ShooterResponse
	Register() error
}
