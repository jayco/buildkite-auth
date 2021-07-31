# buildkite-auth [![Build status](https://badge.buildkite.com/1d777cdb62388e04d43d2e1c2dd821674c8ff4c0f2f6668334.svg?branch=main)](https://buildkite.com/jayco/buildkite-auth)[![Go Reference](https://pkg.go.dev/badge/github.com/jayco/buildkite-auth/v2.svg)](https://pkg.go.dev/github.com/jayco/buildkite-auth/v2)

Simple client to authenticate against the Buildkite API (V2)

# Installation

To get the package:

```shell
go get github.com/jayco/buildkite-auth/v2
```

# API

## Data Structures

### TokenResponse

```go
type TokenResponse struct {
	UUID    string   `json:"uuid,omitempty"`
	Scopes  []string `json:"scopes,omitempty"`
	Message string   `json:"message,omitempty"`
}
```

## Methods

### NewClient(apiToken string, debug bool) *Client

Returns a new Auth Client.

```go
package main

import (
    "log"
    "github.com/jayco/buildkite-auth/v2"
)

func main() {
	c := client.NewClient("your-api-token", false)
}
```

### TokenScopes() (*TokenResponse, error)

Maps to [GET https://buildkite.com/docs/rest-api/access-token](https://buildkite.com/docs/rest-api/access-token)

Get the current token
Returns details about the API Access Token that was used to authenticate the request.

```go
    resp, err := c.TokenScopes()
    log.Println(resp)
```

```shell
➜  go run main.go
2021/07/29 18:25:52 &{12346789-o98u-hy65-okj7-9876yt54re32 [read_agents write_agents read_teams read_artifacts write_artifacts read_builds write_builds read_job_env read_build_logs write_build_logs read_organizations read_pipelines write_pipelines read_user graphql] }
```

### RevokeToken() (*TokenResponse, error)

Maps to [DELETE https://buildkite.com/docs/rest-api/access-token](https://buildkite.com/docs/rest-api/access-token)

Revokes the API Access Token that was used to authenticate the request. Once revoked, the token can no longer be used for further requests.

```go
    resp, err := c.TokenScopes()
    log.Println(resp)
```

```shell
➜  bkstats go run main.go
2021/07/29 19:40:07 &{ [] 204 No Content}
```
