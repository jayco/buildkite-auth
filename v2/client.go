package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
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

// AuthRoundTripper to override default roundtripper https://github.com/golang/go/blob/master/src/net/http/client.go#L62
type authRoundTripper struct {
	Token     string
	Dump      bool
	Transport http.RoundTripper
}

// RoundTrip overrides the default RoundTrip https://github.com/golang/go/blob/master/src/net/http/client.go#L142
func (a authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))

	if a.Dump {
		if dump, err := httputil.DumpRequest(req, true); err == nil {
			fmt.Printf("REQUEST_DUMP:\n%s\n%s\n", req.URL, dump)
		}
	}

	return a.Transport.RoundTrip(req)
}

// NewClient with our roundtripper
func NewClient(apiToken string, debug bool) *Client {
	c := http.Client{Transport: &authRoundTripper{Token: apiToken, Dump: debug, Transport: http.DefaultTransport}}
	return &Client{Client: &c, Host: v2Api}
}
