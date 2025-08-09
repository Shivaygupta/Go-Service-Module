package handlers

import (
	"net/http"
	"strconv"

	"kong-go-assignment/errors"
	"kong-go-assignment/models"
	"kong-go-assignment/services"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	svc services.ServiceService
}

func NewServiceHandler(svc services.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

func (h *ServiceHandler) GetServices(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.Error(errors.ErrInvalidInput)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.Error(errors.ErrInvalidInput)
		return
	}

	filter := models.ServiceListFilter{
		Name:   c.Query("name"),
		SortBy: c.DefaultQuery("sort_by", "created_at"),
		Order:  c.DefaultQuery("order", "desc"),
		Page:   page,
		Limit:  limit,
	}

	filter.Normalize()

	data, err := h.svc.ListServices(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": data})
}

func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.Error(errors.ErrInvalidID)
		return
	}

	data, err := h.svc.GetService(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	c.JSON(200, gin.H{"data": data})
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var svc models.Service
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	if err := h.svc.CreateService(c.Request.Context(), &svc); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": svc})
}
