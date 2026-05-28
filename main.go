package main

import (
	"log"

	"github.com/PriyankaGupta1705/go-url-shortener/handler"
	"github.com/PriyankaGupta1705/go-url-shortener/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file (only needed locally)
	godotenv.Load()

	// connect postgres
	pg, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	log.Println("connected to postgres")

	// run migrations
	if err := pg.RunMigrations(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations complete")

	// connect redis
	rdb, err := store.NewRedisStore()
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	log.Println("connected to redis")

	// setup routes
	h := handler.NewHandler(pg, rdb)
	r := gin.Default()

	r.GET("/health", h.Health)
	r.POST("/shorten", h.Shorten)
	r.GET("/:code", h.Redirect)
	r.GET("/stats/:code", h.Stats) // new endpoint

	log.Println("server running on :8080")
	r.Run(":8080")
}
