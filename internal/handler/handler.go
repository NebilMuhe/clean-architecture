package handler

import (
	"github.com/gin-gonic/gin"
)

type User interface {
	CreateUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}
