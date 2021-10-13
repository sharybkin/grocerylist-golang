package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string, component string) {
	log.Error(component, message)

	if statusCode == 500 {
		message = "internal server error"
	}

	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
