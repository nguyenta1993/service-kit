package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/nguyenta1993/service-kit/http/request"
	"github.com/nguyenta1993/service-kit/logger"
	jwtToken "github.com/nguyenta1993/service-kit/token"
	"go.uber.org/zap"
)

type UnauthorizedErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func NewUnauthorizedErrorResponse() *UnauthorizedErrorResponse {
	return &UnauthorizedErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid token",
	}
}

func UserContextMiddleware(logger logger.Logger) gin.HandlerFunc {
	tokenHeaderName := "Bearer "
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if !strings.Contains(authToken, tokenHeaderName) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewUnauthorizedErrorResponse())
			return
		}
		token := authToken[len(tokenHeaderName):]
		claims, err := jwtToken.ParseTokenUnverify(token)
		if err != nil {
			logger.Error("Parse token error", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewUnauthorizedErrorResponse())
			return
		}

		var userContext request.UserContext
		err = mapstructure.WeakDecode(claims, &userContext)
		if err != nil {
			logger.Error("Token not valid", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewUnauthorizedErrorResponse())
			return
		}
		request.SetUserContext(c, &userContext)
		c.Next()
	}
}
