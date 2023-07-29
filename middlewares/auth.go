package middlewares

import (
	"raychat/settings"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	if settings.Get().ExternalToken == "" {
		c.Next()
		return
	}

	rawtoken := c.GetHeader("Authorization")
	tokenStrlist := strings.Split(rawtoken, " ")
	if len(tokenStrlist) != 2 || len(rawtoken) == 0 {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	token := tokenStrlist[1]
	if token != settings.Get().ExternalToken {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	c.Next()
}
