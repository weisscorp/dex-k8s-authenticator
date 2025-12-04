package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderToken(t *testing.T) {
	cluster := Cluster{
		Name:              "test-cluster",
		Client_ID:         "test-client",
		Client_Secret:     "test-secret",
		K8s_Master_URI:    "https://k8s.example.com:6443",
		K8s_Ca_Pem:       "test-ca-pem",
		Redirect_URI:      "http://localhost:5555/callback",
		Static_Context_Name: false,
		Namespace:         "default",
	}

	claims := map[string]interface{}{
		"iss":   "https://dex.example.com",
		"email": "test@example.com",
		"sub":   "test-user",
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("Failed to marshal claims: %v", err)
	}

	w := httptest.NewRecorder()
	cluster.renderToken(w, "test-id-token", "test-refresh-token", "", "", "", "/", "v1.23.0", claimsJSON)

	if w.Code != http.StatusOK {
		t.Errorf("renderToken() status = %v, want %v", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test-cluster") {
		t.Error("renderToken() body should contain cluster name")
	}
	if !strings.Contains(body, "test@example.com") {
		t.Error("renderToken() body should contain email")
	}
}

func TestRenderTokenWithEmail(t *testing.T) {
	cluster := Cluster{
		Name: "test-cluster",
	}

	claims := map[string]interface{}{
		"iss":   "https://dex.example.com",
		"email": "user@example.com",
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("Failed to marshal claims: %v", err)
	}

	w := httptest.NewRecorder()
	cluster.renderToken(w, "token", "refresh", "", "", "", "/", "v1.23.0", claimsJSON)

	body := w.Body.String()
	// Username should be extracted from email (part before @)
	if !strings.Contains(body, "user") {
		t.Error("renderToken() should extract username from email")
	}
}

func TestRenderTokenWithoutEmail(t *testing.T) {
	cluster := Cluster{
		Name: "test-cluster",
	}

	claims := map[string]interface{}{
		"iss": "https://dex.example.com",
		"sub": "test-user",
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("Failed to marshal claims: %v", err)
	}

	w := httptest.NewRecorder()
	cluster.renderToken(w, "token", "refresh", "", "", "", "/", "v1.23.0", claimsJSON)

	body := w.Body.String()
	// Should use default username "user" when email is not present
	if !strings.Contains(body, "user") {
		t.Error("renderToken() should use default username when email is missing")
	}
}
