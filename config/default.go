package config

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitSession(r *gin.Engine) {
	store := cookie.NewStore([]byte(Config.AppSecret))
	opts := sessions.Options{
		Path:     "/",
		MaxAge:   1800, // 30 Minutes
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if !Config.AppProd {
		opts = sessions.Options{
			Path:     "/",
			MaxAge:   1800, // 30 Minutes
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}
	}

	store.Options(opts)
	r.Use(sessions.Sessions("tz-sessions", store))
}

func SetCORS(r *gin.Engine) {
	setConfig := cors.DefaultConfig()
	setConfig.AllowOrigins = split(Config.AllowOrigins)
	setConfig.AllowHeaders = split(Config.AllowHeaders)
	r.Use(cors.New(setConfig))
}

func split(s string) []string {
	return strings.Split(s, "|")
}
