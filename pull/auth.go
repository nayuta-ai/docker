package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type token struct {
	Token string `json:"token"`
}

var (
	authURL = "https://auth.docker.io/token"
	service = "registry.docker.io"
)

// Auth gets a bearer token from the repository
// using image name. user name, and user password.
func (c *Client) Auth() (string, error) {
	u, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("service", service)
	q.Set("scope", "repository:"+c.Name+":pull,push")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("could not parse authentication url: %w", err)
	}

	// If authentication is needed for your repository,
	// set an user name and password to a request.
	if c.User != "" && c.Password != "" {
		req.SetBasicAuth(c.User, c.Password)
	} else if c.User != "" {
		return "", errors.New("error: must specify your password")
	} else if c.Password != "" {
		return "", errors.New("error: must specify your user name")
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed sending auth request: %w", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %s", rsp.Status)
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read auth response body: %w", err)
	}
	var tok token
	if err := json.Unmarshal(body, &tok); err != nil {
		return "", fmt.Errorf("failed to unmarshal token: %w", err)
	}
	return tok.Token, nil
}
