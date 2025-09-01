package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinLogger_ManyRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinLogger())

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	r.POST("/echo", func(c *gin.Context) {
		var json struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": json.Name, "age": json.Age})
	})

	// 批量 GET
	// for i := 0; i < 50; i++ {
	// 	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	// 	w := httptest.NewRecorder()
	// 	r.ServeHTTP(w, req)
	// 	if w.Code != 200 || w.Body.String() != "pong" {
	// 		t.Errorf("GET /ping failed at %d: code=%d, body=%s", i, w.Code, w.Body.String())
	// 	}
	// }

	// 批量 POST JSON
	body := `{"name": "Alice", "age": 30}`
	for i := 0; i < 50; i++ {
		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 200 || !strings.Contains(w.Body.String(), "Alice") {
			t.Errorf("POST /echo failed at %d: code=%d, body=%s", i, w.Code, w.Body.String())
		}
	}
}
