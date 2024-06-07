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
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type UserTestState struct {
	errorMessage string
	isRegistered bool
}

type User struct {
	Username string
	Email    string
	Password string
}

var server *gin.Engine
var pool *pgxpool.Pool

func IntializeTest() {
	ctx := context.Background()
	log := initiator.InitLogger()
	initiator.InitConfig(ctx, "config", "../../config", log)
	log.Info(ctx, "initializing database")
	pool = initiator.InitDB(ctx, viper.GetString("database.url"), log)
	log.Info(ctx, "initilaizied database")

	log.Info(ctx, "initializing migration")
	initiator.InitMigration(ctx, viper.GetString("database.testPath"), viper.GetString("database.murl"), log)
	log.Info(ctx, "initialized migration")

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
	pool.Exec(context.Background(), "DELETE FROM users;")
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

func (u *UserTestState) iSendRequestToWithPayload(method, url string, body *godog.DocString) error {
	req := httptest.NewRequest(method, url, bytes.NewReader([]byte(body.Content)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	response := model.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		u.errorMessage = response.Error.Message
	} else {
		u.isRegistered = response.OK
	}
	return nil
}

func (u *UserTestState) theResponseShouldBe(err string) error {
	if u.errorMessage == err {
		return nil
	}
	return godog.ErrPending
}

func (u *UserTestState) theSystemReturnABoolean(isRegistered string) error {
	value, _ := strconv.ParseBool(isRegistered)
	if u.isRegistered == value {
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
		ts.reset()
		return ctx, nil
	})

	ctx.Step(`^The system sholud return "([^"]*)"$`, ts.theSystemSholudReturn)
	ctx.Step(`^User enters "([^"]*)",""([^"]*)"", and "([^"]*)"$`, ts.userEntersAnd)
	ctx.Step(`^User is on registration page$`, ts.userIsOnRegistrationPage)

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with payload:$`, ts.iSendRequestToWithPayload)
	ctx.Step(`^the response should be "([^"]*)"$`, ts.theResponseShouldBe)

	ctx.Step(`^the system return a boolean "([^"]*)"$`, ts.theSystemReturnABoolean)

	ctx.Step(`^User enters "([^"]*)" and "([^"]*)"$`, ts.userEntersUsernameAndPassword)
	ctx.Step(`^User is on login page$`, ts.userIsOnLoginPage)
	ctx.Step(`^The system sholud return an error "([^"]*)"$`, ts.theSystemSholudReturnAnError)
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
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
