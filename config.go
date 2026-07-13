package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// config speiler ~/.config/ghx/config.toml:
//
//	[owners]
//	lassevk          = "ghp_..."
//	"larvik-kommune" = "ghp_..."
type config struct {
	Owners map[string]string `toml:"owners"`
}

// configPath returnerer stien til config-fila, med respekt for
// $XDG_CONFIG_HOME.
func configPath() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "ghx", "config.toml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// Fall tilbake til en relativ sti; loadConfig gir en tydelig feil.
		return filepath.Join(".config", "ghx", "config.toml")
	}
	return filepath.Join(home, ".config", "ghx", "config.toml")
}

// loadConfig leser og parser config-fila og returnerer en owner→token-mapping
// der nøklene er normalisert til lowercase (matcher owner-parsingen).
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

	tokens := make(map[string]string, len(c.Owners))
	for owner, token := range c.Owners {
		tokens[strings.ToLower(owner)] = token
	}
	return tokens, nil
}
