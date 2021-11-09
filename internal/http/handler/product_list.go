package handler

import (
	"github.com/sharybkin/grocerylist-golang/internal/model"
	"github.com/sharybkin/grocerylist-golang/pkg/extension"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserLists(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		//Response Body was formed inside getUserId
		return
	}
	listsInfo, err := h.services.UserLists.GetUserLists(userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), "getUserLists")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"productLists": listsInfo,
	})
}

func (h *Handler) updateProductList(c *gin.Context) {

	listId := c.Param("id")

	userId, err := getUserId(c)

	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	var request model.ProductListRequest

	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "updateProductList")
		return
	}

	if err := h.services.ProductList.UpdateProductList(userId, listId, request); err != nil {
		setErrorResponse(c, err, "updateProductList")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) deleteProductList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	listId := c.Param("id")

	if err := h.services.ProductList.DeleteProductList(userId, listId); err != nil {
		setErrorResponse(c, err, "deleteProductList")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) createProductList(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	var productList model.ProductListRequest

	if err := c.BindJSON(&productList); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "createProductList")
		return
	}

	listId, err := h.services.ProductList.CreateProductList(userId, productList)

	if err != nil {
		setErrorResponse(c, err, "createProductList")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"productListId": listId,
	})
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
