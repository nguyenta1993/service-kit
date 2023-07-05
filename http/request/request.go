package request

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	jwtToken "github.com/nguyenta1993/service-kit/token"
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

func MustGetUser(c *gin.Context) (user UserContext) {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		panic("missing token")
	}
	token := authToken[len("Bearer "):]
	claims, err := jwtToken.ParseTokenUnverify(token)
	if err != nil {
		panic("invalid token")
	}
	if err := mapstructure.WeakDecode(claims, &user); err != nil {
		panic("can't decode token")
	}
	return
}
