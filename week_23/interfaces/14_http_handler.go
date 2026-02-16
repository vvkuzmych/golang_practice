package main

import (
	"fmt"
	"net/http"
)

// 14. HTTP Handler Interface - Standard net/http handler pattern

type CustomHandler struct {
	Message string
}

func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s: %s %s", h.Message, r.Method, r.URL.Path)
}

type mockResponseWriter struct {
	body string
}

func (m *mockResponseWriter) Header() http.Header { return http.Header{} }
func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.body = string(b)
	return len(b), nil
}
func (m *mockResponseWriter) WriteHeader(int) {}

func main() {
	handler := &CustomHandler{Message: "Hello"}
	req, _ := http.NewRequest("GET", "/api/items", nil)
	resp := &mockResponseWriter{}
	handler.ServeHTTP(resp, req)
	fmt.Println("Response:", resp.body)
}
