package handler

import (
	"github.com/gin-gonic/gin"
)

type User interface {
	CreateUser(ctx *gin.Context)
}
