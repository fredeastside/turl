package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"turl/internal/url"
)

func TestHealth(t *testing.T) {
	s := NewServer(&url.URLShortener{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("You received a %v error.", w.Code)
	}

	expected := "OK"
	actual := w.Body.String()

	if actual != expected {
		t.Errorf("Response should be %v, was %v.", expected, actual)
	}
}
