package handlers

import "net/http"

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

// @Summary Rate limited route
// @Description Access protected route with rate limiting
// @Tags root
// @Produce plain
// @Param API_KEY header string false "API Key"
// @Success 200 {string} string "Request allowed"
// @Failure 429 {string} string "you have reached the maximum number of requests or actions allowed within a certain time frame"
// @Router / [get]
func (d *DefaultHandler) BaseAccess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Request allowed"))
}
