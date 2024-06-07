package test

import (
	"bytes"
	"clean-architecture/internal/constants/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/cucumber/godog"
)

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
