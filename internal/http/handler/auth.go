package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "SignUp")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {

		if _, ok := err.(*extension.AuthError); ok {

			newErrorResponse(c, http.StatusForbidden, err.Error(), "SignUp")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "SignUp")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "SignIn")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {

		if _, ok := err.(*extension.AuthError); ok {

			newErrorResponse(c, http.StatusUnauthorized, "invalid username or password", "SignIn")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "SignIn")
		return

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
