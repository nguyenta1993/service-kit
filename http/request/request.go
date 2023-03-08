package request

import (
	"github.com/gin-gonic/gin"
)

const UserContextKey = "usercontext"

type UserContext struct {
	Id    int64 `mapstructure:"_id"`
	OrgId int64 `mapstructure:"_orgId"`
}

func GetUserContext(c *gin.Context) UserContext {
	return c.MustGet(UserContextKey).(UserContext)
}

func SetUserContext(c *gin.Context, userContext *UserContext) {
	c.Set(UserContextKey, *userContext)
}
