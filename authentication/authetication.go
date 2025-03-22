package authentication

import (
	"github.com/chuongthanh0410/interview/config"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) bool {
	token := c.GetHeader("Authorization")

	// Kiểm tra token
	if token == "" {
		c.JSON(400,
			gin.H{
				"message": "Token is required",
				"status":  false,
				"data":    nil,
			},
		)
		return false
	}

	if token != config.STATIC_TOKEN {
		c.JSON(401,
			gin.H{
				"message": "Invalid token",
				"status":  false,
				"data":    nil,
			},
		)
		return false
	}

	return true
}
