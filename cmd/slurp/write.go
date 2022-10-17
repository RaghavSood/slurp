package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func writeFile(dir string, path string, uuid string, suffix string, data string) error {
	basePath := filepath.Join(dir, path)
	err := ensurePath(basePath)
	if err != nil {
		return errors.Wrap(err, "could not create directories")
	}

	filePath := filepath.Join(basePath, fmt.Sprintf("%s_%s", uuid, suffix))

	err = os.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}

	return nil
}

func ensurePath(outputDir string) error {
	return os.MkdirAll(outputDir, 0750)
}
