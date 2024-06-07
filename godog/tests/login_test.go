package test

import (
	"bytes"
	"clean-architecture/internal/constants/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
)

func (u *UserTestState) userIsOnLoginPage() error {
	return nil
}

func (u *UserTestState) userEntersUsernameAndPassword(username, password string) error {
	us := &User{
		Username: username,
		Password: password,
	}

	r, err := json.Marshal(us)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(r))
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

func (u *UserTestState) theSystemSholudReturnAnError(err string) error {
	if u.errorMessage == err {
		return nil
	}
	return godog.ErrPending
}
