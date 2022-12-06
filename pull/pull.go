package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/opencontainers/go-digest"
)

const (
	registryUrl = "https://registry.hub.docker.com"
)

func (c *Client) getCurrentManifest(tag string) (*Manifest, error) {
	u := strings.Join([]string{registryUrl, "v2", c.Name, "manifests", tag}, "/")
	req, err := c.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return &Manifest{}, err
	}
	req.Header.Set("Accept", MediaTypeManifest)

	log.Printf("Sending %s", u)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Manifest{}, fmt.Errorf("manifest upload failed: %w", err)
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return &Manifest{}, fmt.Errorf("failed to read body on manifest upload response: %w", err)
	}

	if rsp.StatusCode != http.StatusCreated && rsp.StatusCode != http.StatusOK {
		return &Manifest{}, fmt.Errorf("unexpected status %s. %s", rsp.Status, string(body))
	}
	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return &Manifest{}, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}
	if rsp.StatusCode != http.StatusCreated && rsp.StatusCode != http.StatusOK {
		return &Manifest{}, fmt.Errorf("unexpected status %s. %s", rsp.Status, string(body))
	}
	if err != nil {
		return &Manifest{}, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}
	file, err := os.Create(strings.Join([]string{c.Dir, "manifest.json"}, "/"))
	if err != nil {
		return &Manifest{}, fmt.Errorf("failed to create json file: %w", err)
	}
	defer file.Close()
	file.Write(body)
	return &manifest, nil
}

func (c *Client) getConfig(digest digest.Digest, tag string) (*Image, error) {
	u := strings.Join([]string{registryUrl, "v2", c.Name, "blobs", digest.String()}, "/")
	req, err := c.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return &Image{}, err
	}
	log.Printf("Sending %s", u)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Image{}, fmt.Errorf("manifest upload failed: %w", err)
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return &Image{}, fmt.Errorf("failed to read body on manifest upload response: %w", err)
	}

	if rsp.StatusCode != http.StatusCreated && rsp.StatusCode != http.StatusOK {
		return &Image{}, fmt.Errorf("unexpected status %s. %s", rsp.Status, string(body))
	}
	var config Image
	if err := json.Unmarshal(body, &config); err != nil {
		return &Image{}, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}
	if rsp.StatusCode != http.StatusCreated && rsp.StatusCode != http.StatusOK {
		return &Image{}, fmt.Errorf("unexpected status %s. %s", rsp.Status, string(body))
	}
	file, err := os.Create(strings.Join([]string{c.Dir, string(digest)[7:] + ".json"}, "/"))
	if err != nil {
		return &Image{}, fmt.Errorf("failed to create json file: %w", err)
	}
	defer file.Close()
	file.Write(body)
	return &config, nil
}

func (c *Client) getLayer(digest digest.Digest, tag string) error {
	u := strings.Join([]string{registryUrl, "v2", c.Name, "blobs", digest.String()}, "/")
	req, err := c.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	log.Printf("Sending %s", u)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("layer tar file fetched failed: %w", err)
	}
	if err := os.Mkdir(strings.Join([]string{c.Dir, string(digest)[7:]}, "/"), os.ModePerm); err != nil {
		return fmt.Errorf("failed creating a directory: %w", err)
	}
	out, err := os.Create(strings.Join([]string{c.Dir, string(digest)[7:], "layer.tar.gz"}, "/"))
	if err != nil {
		return fmt.Errorf("failed creating tar file: %w", err)
	}
	defer out.Close()
	defer rsp.Body.Close()
	_, err = io.Copy(out, rsp.Body)
	if err != nil {
		return fmt.Errorf("failed copying tar file: %w", err)
	}
	return nil
}
