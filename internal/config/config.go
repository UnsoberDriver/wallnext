// Package config
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type Config struct {
	Path     string   `toml:"-"`
	Settings Settings `toml:"settings"`
	SWWW     SWWW     `toml:"swww"`
	Hypr     Hypr     `toml:"hyprpaper"`
}

type Settings struct {
	Folder     string   `toml:"folder"`
	Backend    string   `toml:"backend"`
	Mode       string   `toml:"mode"`
	Recursive  bool     `toml:"recursive"`
	Extensions []string `toml:"extensions"`
}

type SWWW struct {
	TransitionType     string  `toml:"transition_type"`
	TransitionDuration float64 `toml:"transition_duration"`
}

type Hypr struct {
	Monitor string `toml:"monitor"`
}

func DefaultPath() string {
	return filepath.Join(xdg.ConfigHome, "wallnext", "config.toml")
}

func ExpandPath(p string) string {
	if p == "" {
		return p
	}
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, _ := os.UserHomeDir()
		if p == "~" {
			return home
		}
		return filepath.Join(home, strings.TrimPrefix(p, "~/"))
	}
	return p
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath()
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if _, err := toml.Decode(string(raw), &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.Path = path
	if err := cfg.applyDefaults(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) applyDefaults() error {
	if c.Settings.Folder == "" {
		home, _ := os.UserHomeDir()
		c.Settings.Folder = filepath.Join(home, "Pictures", "Wallpapers")
	}
	c.Settings.Folder = ExpandPath(c.Settings.Folder)
	if c.Settings.Backend == "" {
		c.Settings.Backend = "auto"
	}
	if c.Settings.Mode == "" {
		c.Settings.Mode = "random"
	}
	if len(c.Settings.Extensions) == 0 {
		c.Settings.Extensions = []string{"jpg", "jpeg", "png", "webp", "gif", "avif"}
	}
	if c.SWWW.TransitionType == "" {
		c.SWWW.TransitionType = "fade"
	}
	if c.SWWW.TransitionDuration == 0 {
		c.SWWW.TransitionDuration = 1.0
	}
	return nil
}
