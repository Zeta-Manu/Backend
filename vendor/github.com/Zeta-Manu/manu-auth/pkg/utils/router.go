package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterWithLogger struct {
	Router *gin.Engine
	Logger *zap.Logger
}
