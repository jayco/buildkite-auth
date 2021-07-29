package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type expected struct {
	ValidToken   TokenResponse
	InvalidToken TokenResponse
}

func apiMock() (*httptest.Server, *expected) {
	validToken := TokenResponse{
		UUID: "12345678-1qas-1wer-2wed-12we34rt56y7",
		Scopes: []string{
			"read_agents", "write_agents", "read_teams", "read_artifacts", "write_artifacts", "read_builds",
			"write_builds", "read_job_env", "read_build_logs", "write_build_logs", "read_organizations",
			"read_pipelines", "write_pipelines", "read_user", "graphql",
		},
	}

	invalidToken := TokenResponse{
		Message: "Authentication required. Please supply a valid API Access Token: https://buildkite.com/docs/apis/rest-api#authentication",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure the token is set correctly for authorization
		if r.Header.Get("Authorization") == "Bearer token" {
			// ensure we are hitting the expected api path
			if r.URL.String() == tokenPath {
				// Get https://buildkite.com/docs/rest-api/access-token
				if r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(validToken)
				}

				// Delete https://buildkite.com/docs/rest-api/access-token
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
				}
			}
		} else {
			// sad path
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(invalidToken)
		}
	}))

	return ts, &expected{validToken, invalidToken}
}

func TestClient_TokenScopes(t *testing.T) {
	ts, expResp := apiMock()
	defer ts.Close()

	tests := []struct {
		name    string
		client  *Client
		want    *TokenResponse
		wantErr bool
	}{
		{
			"Should return scopes with valid token",
			&Client{
				&http.Client{
					Transport: &authRoundTripper{
						Token:     "token",
						Transport: http.DefaultTransport,
					},
				},
				ts.URL,
			},
			&expResp.ValidToken,
			false,
		},
		{
			"Should error scopes with invalid token",
			&Client{
				&http.Client{
					Transport: &authRoundTripper{
						Token:     "invalid-token",
						Transport: http.DefaultTransport,
					},
				},
				ts.URL,
			},
			&expResp.InvalidToken,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.TokenScopes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("Client.GetScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RevokeToken(t *testing.T) {
	ts, _ := apiMock()
	defer ts.Close()

	tests := []struct {
		name    string
		Client  *Client
		want    *TokenResponse
		wantErr bool
	}{
		{
			"Should return success when token is revoked",
			&Client{
				&http.Client{
					Transport: &authRoundTripper{
						Token:     "token",
						Transport: http.DefaultTransport,
					},
				},
				ts.URL,
			},
			&TokenResponse{Message: "204 No Content"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.Client.RevokeToken()

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RevokeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RevokeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	token := "token"
	tests := []struct {
		name  string
		token string
		want  *Client
	}{
		{
			"Should return a new client with correct input",
			token,
			&Client{
				Client: &http.Client{
					Transport: &authRoundTripper{
						Token:     token,
						Transport: http.DefaultTransport,
					},
				},
				Host: v2Api,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.token, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
