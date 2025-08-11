package handlers

import (
	"kong-go-assignment/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *ServiceHandler) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/services", h.GetServices)
		v1.GET("/services/:id", h.GetServiceByID)
		v1.POST("/services", h.CreateService)
		v1.DELETE("/services/:id", h.DeleteServiceByID)
	}

	return r
}
