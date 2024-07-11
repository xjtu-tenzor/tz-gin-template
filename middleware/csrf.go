package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GenerateCSRFToken generates a secure CSRF token
func GenerateCSRFToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// CSRFMiddleware checks for the presence and correctness of a CSRF token
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("CSRF-Token")
		if err != nil || sessionToken == "" {
			newToken := GenerateCSRFToken()
			c.SetCookie("CSRF-Token", newToken, 3600, "/", "", false, true)
			c.Set("CSRF-Token", newToken)
		}

		headerToken := c.GetHeader("X-CSRF-Token")
		if headerToken != sessionToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			return
		}

		c.Next()
	}
}

// SetCSRFToken sets the CSRF token in the user's cookie if not already set or expired
func SetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("CSRF-Token")
		if err != nil {
			token := GenerateCSRFToken()
			c.SetCookie("CSRF-Token", token, 3600, "/", "", false, true)
			c.Set("CSRF-Token", token)
		}
		c.Next()
	}
}
