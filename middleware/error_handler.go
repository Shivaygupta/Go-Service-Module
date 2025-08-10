package middleware

import (
	"log"
	"net/http"

	"kong-go-assignment/errors"

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
			err := c.Errors[0].Err
			log.Printf("Error: %v", err)

			if e, ok := err.(*errors.AppError); ok {
				c.JSON(e.StatusCode, APIResponse{
					Error:      e.Error(),
					Message:    e.Debug,
					StatusCode: e.StatusCode,
				})
			} else {

				c.JSON(http.StatusInternalServerError, APIResponse{
					Error:      err.Error(),
					Message:    errors.ErrInternal.Message,
					StatusCode: http.StatusInternalServerError,
				})
			}
		}
	}
}
