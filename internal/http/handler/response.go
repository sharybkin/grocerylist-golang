package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func setErrorResponse(c *gin.Context, err error, component string) {
	if _, ok := err.(*extension.BadRequestError); ok {

		newErrorResponse(c, http.StatusBadRequest, err.Error(), component)
		return
	}

	if _, ok := err.(*extension.NotFoundError); ok {

		newErrorResponse(c, http.StatusNotFound, err.Error(), component)
		return
	}

	newErrorResponse(c, http.StatusInternalServerError, err.Error(), component)
	return
}
