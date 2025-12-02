package middleware

import (
	"AppointmentAPI/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next() 

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.HTTP, gin.H{
					"error": appErr.Code,
					"msg":   appErr.Message,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal_error",
				"msg":   "something went wrong",
			})
		}
	}
}
