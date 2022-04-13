package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	RegisterEndpoint(*gin.RouterGroup)
}
