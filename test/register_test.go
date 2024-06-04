package test

import (
	"clean-architecture/initiator"
	"context"

	"github.com/cucumber/godog"
)

func IntializeTest() {
	initiator.Initialize(context.Background())
}

func userIsOnRegistrationPage() error {
	return nil
}

func userEntersAnd(username, email, password string) error {
	return godog.ErrPending
}

func theSystemSholudReturn(err string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^The system sholud return "([^"]*)"$`, theSystemSholudReturn)
	ctx.Step(`^User enters "([^"]*)",""([^"]*)"", and "([^"]*)"$`, userEntersAnd)
	ctx.Step(`^User is on registre page$`, userIsOnRegistrationPage)
}
