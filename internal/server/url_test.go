package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockURLShortener struct {
}

func (m *mockURLShortener) Encode(s string) (string, error) {
	return "short", nil
}

func (m *mockURLShortener) Decode(s string) (string, error) {
	return "https://google.com", nil
}
func (m *mockURLShortener) GetDailyCount(s string) (int, error) {
	return 1, nil
}
func (m *mockURLShortener) GetWeeklyCount(s string) (int, error) {
	return 2, nil
}
func (m *mockURLShortener) GetCount(s string) (int, error) {
	return 3, nil
}

func TestMakeShort(t *testing.T) {
	t.Run("Invalid url", func(t *testing.T) {
		s := NewServer(&mockURLShortener{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/", strings.NewReader("url=2"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		s.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("You received a %v error.", w.Code)
		}
	})
	t.Run("Valid url", func(t *testing.T) {
		s := NewServer(&mockURLShortener{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/", strings.NewReader("url=https://google.com"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		s.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("You received a %v error.", w.Code)
		}

		expected := "short"
		actual := w.Body.String()

		if actual != expected {
			t.Errorf("Response should be %v, was %v.", expected, actual)
		}
	})
}

func TestGetLong(t *testing.T) {
	s := NewServer(&mockURLShortener{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/url", nil)
	s.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("You received a %v error.", w.Code)
	}
}

func TestReport(t *testing.T) {
	s := NewServer(&mockURLShortener{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/url/report", nil)
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("You received a %v error.", w.Code)
	}

	expected := "1\n2\n3"
	actual := w.Body.String()

	if actual != expected {
		t.Errorf("Response should be %v, was %v.", expected, actual)
	}
}
