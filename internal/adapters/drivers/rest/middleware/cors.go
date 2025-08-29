// Package middleware
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	origins := []string{
		"http://localhost",
		"http://localhost:8000",
		"http://localhost:8080",
		"http://0.0.0.0:8080",
	}

	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	})
}
