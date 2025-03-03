package io

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
)

// SaveGob encodes the provided data into a gob file stored in the platform-specific data directory.
func SaveGob(filename string, data any) error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(dataDir, filename)

	log.Printf("Saving gob file: %s", fullPath)

	// Create or truncate the file.
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// LoadGob opens the gob file from the platform-specific data directory and decodes its contents into data.
func LoadGob(filename string, data any) error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(dataDir, filename)

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}
