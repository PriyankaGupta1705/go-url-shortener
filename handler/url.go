package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/PriyankaGupta1705/go-url-shortener/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	pg    *store.PostgresStore
	redis *store.RedisStore
}

func NewHandler(pg *store.PostgresStore, redis *store.RedisStore) *Handler {
	return &Handler{pg: pg, redis: redis}
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

// POST /shorten
func (h *Handler) Shorten(c *gin.Context) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	code := generateCode()

	// save to postgres
	if err := h.pg.Save(code, body.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save url"})
		return
	}

	// warm the redis cache immediately
	h.redis.Set(c.Request.Context(), code, body.URL)

	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"short_url":  "http://localhost:8080/" + code,
	})
}

// GET /:code
func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")
	ctx := c.Request.Context()

	// 1. check redis cache first
	if url, hit := h.redis.Get(ctx, code); hit {
		c.Header("X-Cache", "HIT")
		go h.pg.IncrementVisits(code) // async visit count
		c.Redirect(http.StatusMovedPermanently, url)
		return
	}

	// 2. cache miss — go to postgres
	url, err := h.pg.Get(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}

	// 3. store in redis for next time
	h.redis.Set(ctx, code, url)
	go h.pg.IncrementVisits(code)

	c.Header("X-Cache", "MISS")
	c.Redirect(http.StatusMovedPermanently, url)
}

// GET /stats/:code
func (h *Handler) Stats(c *gin.Context) {
	code := c.Param("code")
	original, visits, err := h.pg.GetStats(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":         code,
		"original_url": original,
		"visits":       visits,
		"short_url":    "http://localhost:8080/" + code,
	})
}

// GET /health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
