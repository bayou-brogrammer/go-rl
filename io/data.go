package io

import (
	"os"
	"path/filepath"
)

// DataDirName is the subdirectory where we store our game data.
const DataDirName = "go-rl"

// GetDataDir determines the platform-specific configuration directory,
// appends our game subdirectory, and ensures that it exists.
func GetDataDir() (string, error) {
	// Retrieve the user's config directory (e.g. %AppData% on Windows or ~/.config on Linux/Mac)
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(userConfigDir, DataDirName)
	// Ensure the directory exists; create it if not.
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}
