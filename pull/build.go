package main

import (
	_ "crypto/sha256"
	"fmt"

	"github.com/opencontainers/go-digest"
)

func (c *Client) BuildNewImage() error {
	var tag = c.Tags[0]
	manifest, err := c.getCurrentManifest(tag)
	if err != nil {
		return fmt.Errorf("could not fetch manifest for tag %s: %w", tag, err)
	}
	configDigest := getDigestFromConfig(*manifest)
	_, err = c.getConfig(configDigest, tag)
	if err != nil {
		return fmt.Errorf("could not fetch config for tag %s: %w", tag, err)
	}
	for _, layer := range manifest.Layers {
		layerDigest := layer.Digest
		c.getLayer(layerDigest, tag)
	}
	return nil
}

func getDigestFromConfig(manifest Manifest) digest.Digest {
	return manifest.Config.Digest
}
