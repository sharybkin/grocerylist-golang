package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllProductLists(c *gin.Context) {
	id, err := getUserId(c)

	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	lists, err := h.services.ProductList.GetAllProductLists(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", "getAllProductLists")
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
	})

}

func (h *Handler) getProductListsInfo(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		//Response Body was formed inside getUserId
		return
	}
	listsInfo, err := h.services.UserLists.GetUserLists(userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", "getProductListsInfo")
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"productLists": listsInfo,
	})
}

func (h *Handler) getProductListById(c *gin.Context) {

}

func (h *Handler) updateProductList(c *gin.Context) {

}

func (h *Handler) deleteProductList(c *gin.Context) {

}

func (h *Handler) createProductList(c *gin.Context) {

}
