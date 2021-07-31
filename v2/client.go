package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	// tokenPath of the URL
	tokenPath = "/access-token"

	// v2Api URL
	v2Api = "https://api.buildkite.com/v2"
)

// TokenResponse from API
type TokenResponse struct {
	UUID    string   `json:"uuid,omitempty"`
	Scopes  []string `json:"scopes,omitempty"`
	Message string   `json:"message,omitempty"`
}

// Client for authing
type Client struct {
	*http.Client
	Host string
}

// TokenScopes of the client https://buildkite.com/docs/rest-api/access-token
func (c *Client) TokenScopes() (*TokenResponse, error) {
	resp, err := c.Get(fmt.Sprintf("%s%s", c.Host, tokenPath))
	if err != nil {
		return nil, err
	}

	var response TokenResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// RevokeToken of the client https://buildkite.com/docs/rest-api/access-token
func (c *Client) RevokeToken() (*TokenResponse, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", c.Host, tokenPath), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{Message: resp.Status}, nil
}

// NewClient with bearer token
func NewClient(ctx context.Context, apiToken string) *Client {
	at := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	c := oauth2.NewClient(ctx, at)
	return &Client{Client: c, Host: v2Api}
}
