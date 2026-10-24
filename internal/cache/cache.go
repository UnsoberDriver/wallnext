// Package cache
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type State struct {
	Index     int      `json:"index"`
	History   []string `json:"history"`
	Last      string   `json:"last"`
	All       []string `json:"all"`
}

type Cache struct {
	path string
}

func New() *Cache {
	p := filepath.Join(xdg.StateHome, "wallnext", "state.json")
	return &Cache{path: p}
}

func (c *Cache) Load() (*State, error) {
	raw, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{History: []string{}}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		return &State{History: []string{}}, nil
	}
	if s.History == nil {
		s.History = []string{}
	}
	return &s, nil
}

func (c *Cache) Save(s *State) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	raw, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(c.path, raw, 0o644)
}
