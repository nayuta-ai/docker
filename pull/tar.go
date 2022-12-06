package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func TarDirectory(dirPath string, w io.Writer) error {
	tarWriter := tar.NewWriter(w)
	defer tarWriter.Close()

	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to open directory: %w", err)
	}
	for _, file := range dir {
		if file.IsDir() {
			return fmt.Errorf("directories (%s) are not currently supported", file.Name())
		}
		fi, err := file.Info()
		if err != nil {
			return fmt.Errorf("failed fetching file information: %w", err)
		}
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return fmt.Errorf("failed building tar header for %s: %w", fi.Name(), err)
		}
		header.Name = path.Join("." + header.Name)
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed writing header for %s: %w", fi.Name(), err)
		}

		f, err := os.Open(filepath.Join(dirPath, fi.Name()))
		if err != nil {
			return fmt.Errorf("failed to open directory: %w", err)
		}
		if _, err := io.Copy(tarWriter, f); err != nil {
			f.Close()
			return fmt.Errorf("failed copying %s into tar file: %w", fi.Name(), err)
		}
		f.Close()
	}
	return nil
}
