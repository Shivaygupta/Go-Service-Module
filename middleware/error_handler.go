package middleware

import (
	"kong-go-assignment/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data       any    `json:"data,omitempty"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Error: %v", err)
				if e, ok := err.Err.(*errors.AppError); ok {
					c.JSON(e.StatusCode, APIResponse{
						Error:      e.Error(),
						Message:    e.Debug,
						StatusCode: e.StatusCode,
					})
					return
				}
			}
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:      errors.ErrInternal.Message,
			Message:    errors.ErrInternal.Debug,
			StatusCode: errors.ErrInternal.StatusCode,
		})
	}
}
