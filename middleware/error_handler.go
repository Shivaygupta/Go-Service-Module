package middleware

import (
	"log"
	"net/http"

	"kong-go-assignment/errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			log.Printf("Error: %v", err)

			if e, ok := err.(*errors.AppError); ok {
				c.JSON(e.StatusCode, APIResponse{
					Error:   e.Message,
					Message: e.Debug,
				})
			} else {

				c.JSON(http.StatusInternalServerError, APIResponse{
					Error:   errors.ErrInternal.Message,
					Message: err.Error(),
				})
			}
		}
	}
}
