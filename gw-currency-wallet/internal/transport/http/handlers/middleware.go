package handlers

import (
	"gw-currency-wallet/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(Authorization)
		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, Error{Message: "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, Error{Message: "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := jwt.CheckValidTokenAndGetUserID(token, h.cfg.Secret.Key)
		if err != nil {
			h.logger.Error("invalid token:%v", err)
			c.JSON(http.StatusUnauthorized, Error{Message: "invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)

		c.Next()
	}

}
