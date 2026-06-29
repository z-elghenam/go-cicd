package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey_Valid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey abc123")
	
	apiKey, err := GetAPIKey(headers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiKey != "abc123" {
		t.Errorf("Expected apiKey 'abc123', got '%s'", apiKey)
	}
}

func TestGetAPIKey_NoAuthorizationHeader(t *testing.T) {
	headers := http.Header{}
	
	apiKey, err := GetAPIKey(headers)
	if err == nil {
		t.Error("Expected error for missing Authorization header, got nil")
	}
	if apiKey != "" {
		t.Errorf("Expected empty apiKey, got '%s'", apiKey)
	}
}

func TestGetAPIKey_InvalidAuthFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer token123")
	
	apiKey, err := GetAPIKey(headers)
	if err == nil {
		t.Error("Expected error for invalid auth format, got nil")
	}
	if apiKey != "" {
		t.Errorf("Expected empty apiKey, got '%s'", apiKey)
	}
}

func TestGetAPIKey_EmptyApiKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey ")
	
	apiKey, err := GetAPIKey(headers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiKey != "" {
		t.Errorf("Expected empty apiKey, got '%s'", apiKey)
	}
}
