package httpx

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEnsureHTTPS(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	regularRequest, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	sslRequest, err := http.NewRequest("GET", "https://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	EnsureHTTPS(handler)(w, regularRequest)
	if w.Code == 200 {
		t.Fatalf("Expected failure since scheme was http: got=%d", w.Code)
	}

	w = httptest.NewRecorder()
	EnsureHTTPS(handler)(w, sslRequest)
	if w.Code != 200 {
		t.Fatalf("Expected success since scheme was https: got=%d", w.Code)
	}

	os.Setenv("DISABLE_ENSURE_HTTPS", "1")
	w = httptest.NewRecorder()
	EnsureHTTPS(handler)(w, regularRequest)
	if w.Code != 200 {
		t.Fatalf("Expected success since we've disabled required HTTPS: got=%d", w.Code)
	}
}
