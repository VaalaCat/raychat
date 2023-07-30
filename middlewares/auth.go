package middlewares

import (
	"raychat/settings"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func Auth(c *gin.Context) {
	if len(settings.Get().ExternalToken) == 0 {
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
	if !lo.Contains(settings.Get().ExternalToken, token) {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	c.Next()
}
