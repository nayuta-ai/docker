package main

// NewClient creates a new Client.
func NewClient(o *Options) *Client {
	return &Client{
		Options: *o,
	}
}
