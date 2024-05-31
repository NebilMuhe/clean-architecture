package middleware

import (
	"clean-architecture/internal/constants"
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/joomcode/errorx"
	"github.com/spf13/viper"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			e := ctx.Errors[0]
			err := e.Unwrap()

			constants.ErrorResponse(ctx, CastErrorResponse(err))
			return
		}
	}
}

func ErrorFields(err error) []model.FieldError {
	var errs []model.FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errs = append(errs,
				model.FieldError{
					Name:        i,
					Description: v.Error(),
				},
			)
		}
		return errs
	}
	return nil
}

func CastErrorResponse(err error) *model.ErrorResponse {
	debug := viper.GetBool("debug")
	er := errorx.Cast(err)
	if er == nil {
		return &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "unknown error",
		}
	}

	response := model.ErrorResponse{}
	code, ok := errors.ErrorMap[er.Type()]

	if !ok {
		response = model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "unknown Error",
		}
	} else {
		response = model.ErrorResponse{
			Code:    code,
			Message: er.Message(),
		}
	}

	if debug {
		response.Description = fmt.Sprintf("Error %v", er)
		response.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
	}

	return &response
}
