package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventsHandlerRejectsNonGETRequests(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/events", nil)
	response := httptest.NewRecorder()

	eventsHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if allow := response.Header().Get("Allow"); allow != http.MethodGet {
		t.Fatalf("Allow header = %q, want %q", allow, http.MethodGet)
	}
}

func TestEventsHandlerSetsStreamingHeaders(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/events", nil)
	ctx, cancel := context.WithCancel(request.Context())
	cancel()
	request = request.WithContext(ctx)
	response := httptest.NewRecorder()

	eventsHandler(response, request)

	wants := map[string]string{
		"Content-Type":  "text/event-stream",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
	}
	for name, want := range wants {
		if got := response.Header().Get(name); got != want {
			t.Errorf("%s header = %q, want %q", name, got, want)
		}
	}
}

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	healthHandler(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}
