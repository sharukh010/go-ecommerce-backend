package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/tokens"
)

func Authentication() gin.HandlerFunc{
	return func (c *gin.Context)  {
		ClientTokenArr := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(ClientTokenArr)<2 {
			c.JSON(http.StatusBadRequest,gin.H{"error":"No Authorization Header Provided"})
			c.Abort()
			return 
		}

		ClientToken := ClientTokenArr[1]

		claims,err := tokens.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusBadRequest,gin.H{"error":err})
			c.Abort()
			return 
		}
		c.Set("email",claims.Email)
		c.Set("uid",claims.Uid)
		c.Next()
	}
}