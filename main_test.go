package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewEscrowProvider(t *testing.T) {
	// handler := new(EchoHandler)
	expectedBody := "Hello"

	rw := httptest.NewRecorder()

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://example.com/"),
		strings.NewReader(`{"name": "hokey","address": "crap", "pubkey": "ZnVjaw==", "privkey": "ZnVjaw=="}`),
	)

	if err != nil {
		t.Errorf("Failed to create request.")
	}

	addEscrowProvider(rw, req)

	switch rw.Body.String() {
	case expectedBody:
		// body is equal so no need to do anything
	default:
		t.Errorf("Body (%s) did not match expectation (%s).",
			rw.Body.String(),
			expectedBody)
	}
}
