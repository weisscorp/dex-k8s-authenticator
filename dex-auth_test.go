package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func TestOAuth2Config(t *testing.T) {
	tests := []struct {
		name    string
		cluster Cluster
		want    *oauth2.Config
	}{
		{
			name: "dev mode without provider",
			cluster: Cluster{
				Client_ID:     "test-client",
				Client_Secret: "test-secret",
				Issuer:        "https://dex.example.com",
				Scopes:        []string{"openid", "profile"},
				Redirect_URI:  "http://localhost:5555/callback",
				Provider:      nil,
			},
			want: &oauth2.Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://dex.example.com/auth",
					TokenURL: "https://dex.example.com/token",
				},
				Scopes:      []string{"openid", "profile"},
				RedirectURL: "http://localhost:5555/callback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cluster.oauth2Config()
			if got.ClientID != tt.want.ClientID {
				t.Errorf("oauth2Config().ClientID = %v, want %v", got.ClientID, tt.want.ClientID)
			}
			if got.ClientSecret != tt.want.ClientSecret {
				t.Errorf("oauth2Config().ClientSecret = %v, want %v", got.ClientSecret, tt.want.ClientSecret)
			}
			if got.Endpoint.AuthURL != tt.want.Endpoint.AuthURL {
				t.Errorf("oauth2Config().Endpoint.AuthURL = %v, want %v", got.Endpoint.AuthURL, tt.want.Endpoint.AuthURL)
			}
			if got.Endpoint.TokenURL != tt.want.Endpoint.TokenURL {
				t.Errorf("oauth2Config().Endpoint.TokenURL = %v, want %v", got.Endpoint.TokenURL, tt.want.Endpoint.TokenURL)
			}
			if got.RedirectURL != tt.want.RedirectURL {
				t.Errorf("oauth2Config().RedirectURL = %v, want %v", got.RedirectURL, tt.want.RedirectURL)
			}
		})
	}
}

func TestHandleIndex(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		requestPath    string
		expectedStatus int
		expectedRedirect bool
	}{
		{
			name: "single cluster redirect",
			config: Config{
				Web_Path_Prefix: "/",
				Clusters: []Cluster{
					{Name: "test-cluster"},
				},
			},
			requestPath:     "/",
			expectedStatus:  http.StatusSeeOther,
			expectedRedirect: true,
		},
		{
			name: "multiple clusters render index",
			config: Config{
				Web_Path_Prefix: "/",
				Clusters: []Cluster{
					{Name: "cluster1"},
					{Name: "cluster2"},
				},
			},
			requestPath:     "/",
			expectedStatus:  http.StatusOK,
			expectedRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			tt.config.handleIndex(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleIndex() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedRedirect {
				location := w.Header().Get("Location")
				if location == "" {
					t.Error("handleIndex() expected redirect but Location header is empty")
				}
			}
		})
	}
}

func TestHandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		cluster        Cluster
		devMode        bool
		expectedStatus int
		expectRedirect bool
	}{
		{
			name: "dev mode redirect",
			cluster: Cluster{
				Name:         "test-cluster",
				Issuer:       "https://dex.example.com",
				Client_ID:    "test-client",
				Redirect_URI: "http://localhost:5555/callback",
				Provider:     nil,
			},
			devMode:        true,
			expectedStatus: http.StatusSeeOther,
			expectRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/login/test-cluster", nil)
			w := httptest.NewRecorder()

			tt.cluster.handleLogin(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleLogin() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectRedirect {
				location := w.Header().Get("Location")
				if location == "" {
					t.Error("handleLogin() expected redirect but Location header is empty")
				}
			}
		})
	}
}

func TestHandleCallbackDevMode(t *testing.T) {
	cluster := Cluster{
		Name:         "test-cluster",
		Issuer:       "https://dex.example.com",
		Client_ID:    "test-client",
		Redirect_URI: "http://localhost:5555/callback",
		Provider:     nil,
		Config: Config{
			IDP_Ca_Pem:      "",
			IDP_Ca_Pem_File: "",
			IDP_Ca_URI:      "",
			Logo_Uri:        "",
			Web_Path_Prefix: "/",
			Kubectl_Version: "v1.23.0",
		},
	}

	tests := []struct {
		name           string
		method         string
		code           string
		expectedStatus int
	}{
		{
			name:           "valid mock code",
			method:         "GET",
			code:           "mock-dev-code",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid code",
			method:         "GET",
			code:           "invalid-code",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing code",
			method:         "GET",
			code:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/callback?code="+tt.code+"&state="+exampleAppState, nil)
			w := httptest.NewRecorder()

			cluster.handleCallback(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleCallback() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestRenderHTMLError(t *testing.T) {
	cluster := Cluster{
		Name: "test-cluster",
		Config: Config{
			Logo_Uri:        "http://example.com/logo.png",
			Web_Path_Prefix: "/",
		},
	}

	w := httptest.NewRecorder()

	cluster.renderHTMLError(w, "Test error", http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("renderHTMLError() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("renderHTMLError() Content-Type = %v, want text/html; charset=utf-8", contentType)
	}

	xContentType := w.Header().Get("X-Content-Type-Options")
	if xContentType != "nosniff" {
		t.Errorf("renderHTMLError() X-Content-Type-Options = %v, want nosniff", xContentType)
	}
}

// Test helper to create a test context with HTTP client
func createTestContext(client *http.Client) context.Context {
	return context.WithValue(context.Background(), oauth2.HTTPClient, client)
}
