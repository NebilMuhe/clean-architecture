package helpers

import (
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/utils/logger"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx context.Context, password string, log logger.Logger) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error(ctx, "unable to hash password", zap.Error(err))
		err := errors.ErrWriteError.Wrap(err, "unable to hash password")
		return "", err
	}
	return string(bytes), err
}

func CheckPassword(ctx context.Context, hash, providedPassword string, log logger.Logger) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(providedPassword))
	if err != nil {
		log.Error(ctx, "invalid password", zap.Error(err))
		errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		return err
	}
	return nil
}

func CreateToken(ctx context.Context, id, username string, log logger.Logger) (*usermodel.Token, error) {
	accessToken, err := GenerateAccessToken(id, username)
	if err != nil {
		log.Error(ctx, "unable to create token", zap.Error(err))
		err := errors.ErrUnableToCreate.Wrap(err, "unable to create token")
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(id, username)
	if err != nil {
		log.Error(ctx, "unable to create token", zap.Error(err))
		err := errors.ErrUnableToCreate.Wrap(err, "unable to create token")
		return nil, err
	}

	return &usermodel.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

func VerifyToken(ctx context.Context, tokenString string, log logger.Logger) error {
	secretKey := []byte(viper.GetString("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Error(ctx, "invalid token", zap.Error(err))
		err := errors.ErrBadRequest.Wrap(err, "bad request")
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		err := errors.ErrBadRequest.New("invalid token")
		log.Error(ctx, "invalid token", zap.Error(err))
		return err
	}

	return nil
}

func ExtractUsernameAndID(ctx context.Context, tokenString string, log logger.Logger) (map[string]string, error) {
	secretKey := []byte(viper.GetString("secret_key"))
	var err error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}

		return secretKey, nil
	})

	if err != nil {
		log.Error(ctx, "invalid token", zap.Error(err))
		err := errors.ErrBadRequest.Wrap(err, "bad request")
		return nil, err
	}

	var username string
	var id string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["sub"]; ok {
			if int(claims["sub"].(float64)) == 1 {
				id = claims["id"].(string)
				username = claims["username"].(string)
			}
		} else {
			err := errors.ErrBadRequest.New("invalid token")
			log.Error(ctx, "invalid token", zap.Error(err))
			return nil, err
		}
	}

	return map[string]string{"id": id, "username": username}, nil
}

func Encrypt(ctx context.Context, key []byte, token string, log logger.Logger) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(ctx, "unable to create cipher", zap.Error(err))
		err := errors.ErrUnableToCreate.Wrap(err, "unable to create cipher")
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(token))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Error(ctx, "unable to read", zap.Error(err))
		err := errors.ErrReadError.Wrap(err, "unable to read")
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(token))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ctx context.Context, key []byte, token string, log logger.Logger) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		log.Error(ctx, "unable to decode", zap.Error(err))
		err := errors.ErrReadError.Wrap(err, "unable to read")
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(ctx, "unable to create cipher", zap.Error(err))
		err := errors.ErrUnableToCreate.Wrap(err, "unable to create cipher")
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		err := errors.ErrBadRequest.New("bad request")
		log.Error(ctx, "cipher text less than block size", zap.Error(err))
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
