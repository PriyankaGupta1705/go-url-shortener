package main

import (
	"log"

	"github.com/PriyankaGupta1705/go-url-shortener/handler"
	"github.com/PriyankaGupta1705/go-url-shortener/store"
	"github.com/gin-gonic/gin"
)

func main() {
	s := store.NewStore()
	h := handler.NewHandler(s)

	r := gin.Default()

	r.GET("/health", h.Health)
	r.POST("/shorten", h.Shorten)
	r.GET("/:code", h.Redirect)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
