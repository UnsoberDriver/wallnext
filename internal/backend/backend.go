// Package backend
package backend

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/UnsoberDriver/wallnext/internal/config"
)

type Backend interface {
	Set(path string) error
	Name() string
}

func Detect(c *config.Config) (Backend, error) {
	if c.Settings.Backend != "" && c.Settings.Backend != "auto" {
		return get(c.Settings.Backend, c)
	}
	// try in order
	for _, name := range []string{"swww", "hyprpaper", "kde", "feh", "nitrogen"} {
		if has(name) {
			return get(name, c)
		}
	}
	return nil, fmt.Errorf("no supported wallpaper backend found (install swww, hyprpaper, feh, nitrogen or use KDE)")
}

func has(name string) bool {
	switch name {
	case "swww":
		return commandExists("swww")
	case "hyprpaper":
		return commandExists("hyprpaper")
	case "kde":
		return commandExists("plasma-apply-wallpaperimage") || commandExists("qdbus") || commandExists("qdbus6")
	case "feh":
		return commandExists("feh")
	case "nitrogen":
		return commandExists("nitrogen")
	}
	return false
}

func commandExists(c string) bool {
	_, err := exec.LookPath(c)
	return err == nil
}

func get(name string, c *config.Config) (Backend, error) {
	switch name {
	case "swww":
		return &Swww{cfg: c}, nil
	case "hyprpaper":
		return &Hypr{cfg: c}, nil
	case "kde":
		return &KDE{}, nil
	case "feh":
		return &Feh{}, nil
	case "nitrogen":
		return &Nitrogen{}, nil
	}
	return nil, fmt.Errorf("unknown backend %s", name)
}

type Swww struct {
	cfg *config.Config
}

func (s *Swww) Name() string { return "swww" }
func (s *Swww) Set(path string) error {
	// try to init if not running? swww init is usually done on startup
	args := []string{"img", path, "--transition-type", s.cfg.SWWW.TransitionType, "--transition-duration", fmt.Sprintf("%.2f", s.cfg.SWWW.TransitionDuration)}
	cmd := exec.Command("swww", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("swww: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

type Hypr struct {
	cfg *config.Config
}

func (h *Hypr) Name() string { return "hyprpaper" }
func (h *Hypr) Set(path string) error {
	mon := h.cfg.Hypr.Monitor
	if mon == "" {
		mon = ","
	}
	// hyprpaper is typically controlled via hyprctl
	cmd := exec.Command("hyprctl", "hyprpaper", "wallpaper", mon+","+path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// try unload/reload pattern
		_ = exec.Command("hyprctl", "hyprpaper", "preload", path).Run()
		cmd = exec.Command("hyprctl", "hyprpaper", "wallpaper", mon+","+path)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("hyprpaper: %w (%s)", err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

type KDE struct{}

func (k *KDE) Name() string { return "kde" }
func (k *KDE) Set(path string) error {
	if commandExists("plasma-apply-wallpaperimage") {
		cmd := exec.Command("plasma-apply-wallpaperimage", path)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("plasma-apply-wallpaperimage: %w (%s)", err, strings.TrimSpace(string(out)))
		}
		return nil
	}
	// fallback to dbus for older setups
	if commandExists("qdbus6") {
		// this is a best effort
		cmd := exec.Command("qdbus6", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", fmt.Sprintf("var Desktops = desktops(); for (i=0;i<Desktops.length;i++) { d = Desktops[i]; d.wallpaperPlugin = 'org.kde.image'; d.currentConfigGroup = Array('Wallpaper', 'org.kde.image', 'General'); d.writeConfig('Image', 'file://%s') }", path))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("kde dbus: %w (%s)", err, strings.TrimSpace(string(out)))
		}
		return nil
	}
	return fmt.Errorf("no kde wallpaper tool found")
}

type Feh struct{}

func (f *Feh) Name() string { return "feh" }
func (f *Feh) Set(path string) error {
	cmd := exec.Command("feh", "--bg-fill", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("feh: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

type Nitrogen struct{}

func (n *Nitrogen) Name() string { return "nitrogen" }
func (n *Nitrogen) Set(path string) error {
	cmd := exec.Command("nitrogen", "--set-zoom-fill", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nitrogen: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}
