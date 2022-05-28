package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.POST("/sis", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain; charset=utf-8", []byte("OK"))
	})
	server := &http.Server{
		Handler: router,
	}
	server.ListenAndServe()
}
