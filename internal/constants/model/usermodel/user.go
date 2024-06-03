package usermodel

import (
	"clean-architecture/internal/constants/errors"
	"clean-architecture/utils/logger"
	"context"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RegisterUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginUser struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RefreshToken struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

var usernameRule = []validation.Rule{
	validation.Required.Error("username required"),
	validation.Length(5, 20).Error("Username length must be atleast 5 characters"),
	validation.Match(regexp.MustCompile(`^[A-Za-z]\w{4,}$`)).Error("username must be valid"),
	is.Alphanumeric,
}

var emailRule = []validation.Rule{
	validation.Required.Error("email required"),
	is.Email.Error("email must be valid"),
}

var passwordRule = []validation.Rule{
	validation.Required.Error("password required"),
	validation.Length(8, 50).Error("Password length must be atleast 8 characters long"),
	validation.Match(regexp.MustCompile(`[A-Z]`)).Error("Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters"),
	validation.Match(regexp.MustCompile(`[a-z]`)).Error("Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters"),
	validation.Match(regexp.MustCompile(`[0-9]`)).Error("Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters"),
	validation.Match(regexp.MustCompile(`[-\#\$\.\%\&\*]`)).Error("Password must contain atleast one uppercase letters,one lowercase letters, digits and special characters"),
}

func (u RegisterUser) Validate(ctx context.Context, log logger.Logger) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, usernameRule...),
		validation.Field(&u.Email, emailRule...),
		validation.Field(&u.Password, passwordRule...),
	)

	if err != nil {
		log.Error(ctx, "validation failed", zap.Error(err), zap.String("username", u.Username), zap.String("email", u.Email))
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		return err
	}

	return nil
}

func (u UpdateUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, usernameRule...),
		validation.Field(&u.Email, emailRule...),
		validation.Field(&u.Password, passwordRule...),
	)
}

func (u LoginUser) Validate(ctx context.Context, log logger.Logger) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, usernameRule...),
		validation.Field(&u.Password, passwordRule...))
	if err != nil {
		log.Error(ctx, "validation failed", zap.Error(err), zap.String("username", u.Username))
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		return err
	}

	return nil
}
