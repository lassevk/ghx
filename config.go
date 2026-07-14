package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// config mirrors ~/.config/ghx/config.toml:
//
//	[tokens]
//	lassevk          = "ghp_..."
//	"larvik-kommune" = "ghp_..."
//	"lassevk/ghx"    = "ghp_..."   # per-repo, overrides the owner token
//
// A key without "/" is owner-level; a key with "/" is repo-specific (owner/repo).
type config struct {
	Tokens map[string]string `toml:"tokens"`
}

// configPath returns the path to the config file, respecting $XDG_CONFIG_HOME.
func configPath() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "ghx", "config.toml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// Fall back to a relative path; loadConfig will report a clear error.
		return filepath.Join(".config", "ghx", "config.toml")
	}
	return filepath.Join(home, ".config", "ghx", "config.toml")
}

// loadConfig reads and parses the config file and returns a key->token mapping
// where the keys (owner or owner/repo) are normalized to lowercase (matching
// the owner/repo parsing).
func loadConfig() (map[string]string, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("could not read config %s: %w", path, err)
	}

	var c config
	if err := toml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	tokens := make(map[string]string, len(c.Tokens))
	for key, token := range c.Tokens {
		tokens[strings.ToLower(key)] = token
	}
	return tokens, nil
}
