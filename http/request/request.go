package request

import (
	"github.com/gin-gonic/gin"
)

const UserContextKey = "usercontext"

type UserContext struct {
	UserId      int64  `mapstructure:"sub"`
	Phonenumber string `mapstructure:"phone"`
	IsAnonymous bool
}

func GetUserContext(c *gin.Context) UserContext {
	return c.MustGet(UserContextKey).(UserContext)
}

func SetUserContext(c *gin.Context, userContext *UserContext) {
	if userContext.Phonenumber == "" {
		userContext.IsAnonymous = true
	}
	c.Set(UserContextKey, *userContext)
}
