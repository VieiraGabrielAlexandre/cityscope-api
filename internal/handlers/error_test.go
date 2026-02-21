package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	// Injetar um request ID manual para teste
	req = req.WithContext(SetRequestID(req.Context(), "test-id"))

	rr := httptest.NewRecorder()

	WriteError(rr, req, "TEST_CODE", "test message", http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resp.Error.Code != "TEST_CODE" {
		t.Errorf("wrong error code: got %v want %v", resp.Error.Code, "TEST_CODE")
	}

	if resp.Error.Message != "test message" {
		t.Errorf("wrong error message: got %v want %v", resp.Error.Message, "test message")
	}

	if resp.Error.RequestID != "test-id" {
		t.Errorf("wrong request id: got %v want %v", resp.Error.RequestID, "test-id")
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("wrong content type: got %v want %v", contentType, "application/json; charset=utf-8")
	}
}
