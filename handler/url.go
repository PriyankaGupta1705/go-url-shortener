package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/PriyankaGupta1705/go-url-shortener/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
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
	h.store.Save(code, body.URL)

	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"short_url":  "http://localhost:8080/" + code,
	})
}

// GET /:code
func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")
	url, err := h.store.Get(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}

// GET /health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
