package middlewares

import (
	"github.com/tikivn/s14e-backend-utils/localization"

	"github.com/gin-gonic/gin"
)

func SetLanguage(resources []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Request.FormValue("lang")
		accept := c.Request.Header.Get("Accept-Language")
		localization.NewLocalizer(localization.ResourceConfig{
			Lang:      lang,
			Accept:    accept,
			Resources: resources,
		})
	}
}
