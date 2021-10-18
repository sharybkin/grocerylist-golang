package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sharybkin/grocerylist-golang/internal/service"
)

type Handler struct {
	services *service.ServicesHolder
}

func NewHandler(services *service.ServicesHolder) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/productExamples", h.getAllProductExamples)

		lists := api.Group("/productList")
		{
			lists.GET("/", h.getUserLists)
			lists.GET("/:id", h.getProductListById)
			lists.POST("/", h.createProductList)
			lists.PUT("/:id", h.updateProductList)
			lists.DELETE("/:id", h.deleteProductList)

			products := lists.Group(":id/products")
			{
				products.POST("/", h.createProduct)
				products.GET("/", h.getAllProducts)
				products.PUT("/:product_id", h.updateProduct)
				products.DELETE("/:product_id", h.deleteProduct)
			}

		}
	}

	return router
}
