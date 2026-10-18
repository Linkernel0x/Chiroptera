package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Dependencies contains dependencies' types
type Dependencies struct {
	System []string `toml:"system,omitempty"` // E.g., ["apt:curl", "brew:git"]
	Chiro  []string `toml:"chiro,omitempty"`  // E.g., ["alice@mio-tool#1.0.0"]
}

// Security contains types for security checks
type Security struct {
	SHA256    string `toml:"sha256"`
	Signature string `toml:"signature,omitempty"`
}

// Source contains types for packets' sources
type Source struct {
	Url     string `toml:"url"`
	Torrent string `toml:"torrent,omitempty"`
}

// Recipe defines recipes structure
type Recipe struct {
	Name         string       `toml:"name"`
	Version      string       `toml:"version"`
	Description  string       `toml:"description"`
	Source       Source       `toml:"source"`
	Security     Security     `toml:"security"`
	Dependencies Dependencies `toml:"dependencies,omitempty"`
}

// ParseRecipe reads the TOML end converts it into a Recipe
func ParseRecipe(filePath string) (*Recipe, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	var r Recipe
	if err := toml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("error while parsing TOML: %w", err)
	}

	return &r, nil
}

// WriteRecipeFile writes a recipe to a toml configuration file
func WriteRecipeFile(filePath string, r *Recipe) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error while creating directory: %w", err)
	}

	data, err := toml.Marshal(r)
	if err != nil {
		return fmt.Errorf("error while serializing to TOML: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("error while writing file: %w", err)
	}

	return nil
}
