package test

import (
	"clean-architecture/initiator"
	"context"
	"testing"

	"github.com/cucumber/godog"
)

type UserTestState struct {
}

func IntializeTest() {
	initiator.Initialize(context.Background())
}

func (u *UserTestState) reset() {

}

func (u *UserTestState) userIsOnRegistrationPage() error {
	return nil
}

func (u *UserTestState) userEntersAnd(username, email, password string) error {
	return godog.ErrPending
}

func (u *UserTestState) theSystemSholudReturn(err string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ts := &UserTestState{}
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ts.reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^The system sholud return "([^"]*)"$`, ts.theSystemSholudReturn)
	ctx.Step(`^User enters "([^"]*)",""([^"]*)"", and "([^"]*)"$`, ts.userEntersAnd)
	ctx.Step(`^User is on registre page$`, ts.userIsOnRegistrationPage)
}

func IntializeTestSuite(sc *godog.TestSuiteContext) {
	sc.BeforeSuite(func() {

	})

	sc.AfterSuite(func() {

	})
}

func TestFeautres(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: IntializeTestSuite,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
