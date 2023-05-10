package main

import (
	"bytes"
	"log"
)

func main() {
	o := Options{
		Dir:  "./test",
		Name: "library/busybox",
		Tags: []string{"latest"},
	}
	client := NewClient(&o)
	token, err := client.Auth()
	if err != nil {
		log.Fatalf("failed to authorize. %s", err)
	}
	client.Token = func() string { return token }
	b := &bytes.Buffer{}
	if err := TarDirectory(client.Dir, b); err != nil {
		log.Fatalf("failed to tar layer. %s", err)
	}
	err = client.BuildNewImage()
	if err != nil {
		log.Fatal(err)
	}
}
