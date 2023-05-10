package main

import (
	"io"
	"net/http"
)

// Options contains configuration options for the client
type Options struct {
	// Dir is the directory that we build the container from
	Dir string
	// Name is the name of the repository
	Name string
	// BaseURL is the base URL of the repository. For Docker this is https://index.docker.io
	// For GCR it is https://gcr.io
	BaseURL string
	//
	User     string
	Password string
	// Token is the bearer token for the repository. For GCR you can use $(gcloud auth print-access-token).
	// For Docker, supply your Docker Hub username and password instead.
	Token func() string
	// Tag is the tag for the image. Set to "latest" if you're out of ideas
	Tags []string
}

// Client lets you send a container up to a repository
type Client struct {
	Options
}

// New creates a new Client
func New(o *Options) *Client {
	return &Client{
		Options: *o,
	}
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.Token != nil {
		req.Header.Set("Authorization", "Bearer "+c.Token())
	}

	return req, nil
}
