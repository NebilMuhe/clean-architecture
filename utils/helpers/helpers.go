package helpers

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hash, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func CreateToken(id, username string) (map[string]string, error) {
	accessToken, err := GenerateAccessToken(id, username)
	if err != nil {
		// logger.Error(ctx, "unable to create access token", zap.Error(err))
		// err = errors.ErrUnableToCreate.Wrap(err, "unable to create access token")
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(id, username)
	if err != nil {
		// logger.Error(ctx, "unable to create refresh token", zap.Error(err))
		// err = errors.ErrUnableToCreate.Wrap(err, "unable to create refresh token")
		return nil, err
	}

	return map[string]string{"access_token": accessToken, "refresh_token": refreshToken}, nil
}

func GenerateAccessToken(id, username string) (string, error) {
	secretKey := []byte(viper.GetString("secret_key"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id,
			"username": username,
			"exp":      time.Now().Add(time.Minute * 15).Unix(),
		})

	accessToken, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(id, username string) (string, error) {
	secretKey := []byte(viper.GetString("secret_key"))

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["username"] = username
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add((time.Hour * 24) * 30).Unix()
	rt, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return rt, nil
}
