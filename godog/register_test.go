package test

import (
	"bytes"
	"clean-architecture/initiator"
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/constants/model"
	"clean-architecture/internal/handler/middleware"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cucumber/godog"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserTestState struct {
	errorMessage string
}

type User struct {
	Username string
	Email    string
	Password string
}

var server *gin.Engine

func IntializeTest() {
	ctx := context.Background()
	log := initiator.InitLogger()
	initiator.InitConfig(ctx, "config", "../config", log)
	log.Info(ctx, "initializing database")
	pool := initiator.InitDB(ctx, viper.GetString("database.url"), log)
	log.Info(ctx, "initilaizied database")

	log.Info(ctx, "initializing persistence layer")
	persitence := initiator.InitPersistence(dbinstance.New(pool), log)
	log.Info(ctx, "initialized persistence layer")

	log.Info(ctx, "initializing service layer")
	service := initiator.InitService(persitence, log)
	log.Info(ctx, "initialized service layer")

	log.Info(ctx, "initializing handler layer")
	handler := initiator.InitHandler(service, log)
	log.Info(ctx, "initialized handler")

	log.Info(ctx, "intializing server")
	server = gin.New()
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin-recovery"), true))
	server.Use(middleware.ErrorHandler())
	log.Info(ctx, "initialized server")

	log.Info(ctx, "initializing routes")
	router := server.Group("/api/v1")
	initiator.InitRoute(router, handler)
	log.Info(ctx, "initialized routes")
}

func (u *UserTestState) reset() {

}

func (u *UserTestState) userIsOnRegistrationPage() error {
	return nil
}

func (u *UserTestState) userEntersAnd(username, email, password string) error {
	us := &User{
		Username: username,
		Email:    email,
		Password: password,
	}

	r, err := json.Marshal(us)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewReader(r))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	errorResponse := model.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		return err
	}

	u.errorMessage = errorResponse.Error.FieldError[0].Description

	return nil
}

func (u *UserTestState) theSystemSholudReturn(err string) error {
	if u.errorMessage == err {
		return nil
	}
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
		IntializeTest()
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
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
