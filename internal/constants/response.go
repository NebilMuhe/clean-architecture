package constants

import (
	"clean-architecture/internal/constants/model"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, err *model.ErrorResponse) {
	ctx.AbortWithStatusJSON(err.Code, model.Response{
		OK:    false,
		Error: err,
	})
}

func SuccessResponse(ctx *gin.Context, statusCode int, data interface{}, metadata *model.MetaData) {
	ctx.JSON(statusCode, model.Response{
		OK:       true,
		MetaData: metadata,
		Data:     data,
	})
}
