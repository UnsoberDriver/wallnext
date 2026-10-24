// Package scanner
package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/UnsoberDriver/wallnext/internal/config"
)

func Find(cfg *config.Config) ([]string, error) {
	var out []string
	ext := map[string]bool{}
	for _, e := range cfg.Settings.Extensions {
		ext["."+strings.ToLower(e)] = true
	}

	walk := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if !cfg.Settings.Recursive && p != cfg.Settings.Folder {
				return filepath.SkipDir
			}
			return nil
		}
		if ext[strings.ToLower(filepath.Ext(info.Name()))] {
			out = append(out, p)
		}
		return nil
	}

	if err := filepath.Walk(cfg.Settings.Folder, walk); err != nil {
		return nil, err
	}
	return out, nil
}
