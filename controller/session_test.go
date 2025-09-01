package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	return r
}

func TestSession_ComplexStruct(t *testing.T) {
	r := setupRouter()

	type Profile struct {
		ID    int
		Name  string
		Roles []string
		Meta  map[string]any
	}

	profile := Profile{
		ID:    101,
		Name:  "bob",
		Roles: []string{"admin", "user"},
		Meta:  map[string]any{"age": 30, "active": true},
	}

	r.POST("/set_struct", func(c *gin.Context) {
		SessionSet(c, "profile", profile)
		c.String(200, "ok")
	})

	r.GET("/get_struct", func(c *gin.Context) {
		val := SessionGet(c, "profile")
		if val == nil {
			c.String(404, "not found")
			return
		}
		p, ok := val.(Profile)
		if !ok {
			c.String(500, "type error")
			return
		}
		c.JSON(200, p)
	})

	r.POST("/update_struct", func(c *gin.Context) {
		val := SessionGet(c, "profile")
		if val == nil {
			c.String(404, "not found")
			return
		}
		p := val.(Profile)
		p.Name = "alice"
		p.Meta["age"] = 31
		SessionUpdate(c, "profile", p)
		c.String(200, "updated")
	})

	// set struct
	req := httptest.NewRequest("POST", "/set_struct", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("set struct failed: %d", w.Code)
	}
	cookieVal := w.Header().Get("Set-Cookie")

	// get struct
	req = httptest.NewRequest("GET", "/get_struct", nil)
	req.Header.Set("Cookie", cookieVal)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if newCookie := w.Header().Get("Set-Cookie"); newCookie != "" {
		cookieVal = newCookie
	}
	if !strings.Contains(w.Body.String(), "bob") {
		t.Fatalf("get struct failed: %s", w.Body.String())
	}

	// update struct
	req = httptest.NewRequest("POST", "/update_struct", nil)
	req.Header.Set("Cookie", cookieVal)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if newCookie := w.Header().Get("Set-Cookie"); newCookie != "" {
		cookieVal = newCookie
	}
	if w.Code != 200 {
		t.Fatalf("update struct failed: %d", w.Code)
	}

	// get updated struct
	req = httptest.NewRequest("GET", "/get_struct", nil)
	req.Header.Set("Cookie", cookieVal)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !strings.Contains(w.Body.String(), "alice") || !strings.Contains(w.Body.String(), "31") {
		t.Fatalf("get updated struct failed: %s", w.Body.String())
	}
}
