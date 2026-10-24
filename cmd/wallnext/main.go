// Package main
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/UnsoberDriver/wallnext/internal/backend"
	"github.com/UnsoberDriver/wallnext/internal/cache"
	"github.com/UnsoberDriver/wallnext/internal/config"
	"github.com/UnsoberDriver/wallnext/internal/scanner"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println("wallnext", version)
		return
	case "init":
		printInit()
		return
	case "list":
		runList()
		return
	case "next", "random", "prev", "set":
		runSet(cmd, os.Args[2:])
		return
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`wallnext - press Alt+Shift to go to the next wallpaper

Usage:
  wallnext next    go to next (random or sequential)
  wallnext random  pick a random wallpaper
  wallnext prev    go to previous
  wallnext list    list wallpapers
  wallnext set <path>
  wallnext init    print starter config
  wallnext version`)
}

func printInit() {
	fmt.Println(`# ~/.config/wallnext/config.toml
[settings]
folder = "~/Images/Wallpapers"
backend = "auto"  # auto, swww, kde, hyprpaper, feh, nitrogen
mode = "random"   # random or sequential
recursive = false
extensions = ["jpg", "jpeg", "png", "webp", "gif", "avif"]

[swww]
transition_type = "fade"
transition_duration = 1.0`)
}

func load() (*config.Config, *cache.Cache, *cache.State, []string, backend.Backend, error) {
	cfg, err := config.Load("")
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	c := cache.New()
	st, err := c.Load()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	list, err := scanner.Find(cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if len(list) == 0 {
		return nil, nil, nil, nil, nil, fmt.Errorf("no wallpapers found in %s", cfg.Settings.Folder)
	}
	st.All = list
	b, err := backend.Detect(cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return cfg, c, st, list, b, nil
}

func runList() {
	_, _, _, list, _, err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, p := range list {
		fmt.Println(p)
	}
}

func runSet(cmd string, args []string) {
	cfg, c, st, list, b, err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	var target string
	if cmd == "set" && len(args) == 1 {
		target = args[0]
	} else if cmd == "random" || (cmd == "next" && cfg.Settings.Mode == "random") {
		// pick different from last if possible
		target = list[rand.Intn(len(list))]
		if len(list) > 1 && target == st.Last {
			target = list[rand.Intn(len(list))]
		}
	} else if cmd == "prev" && len(st.History) > 1 {
		// go back one
		h := st.History
		target = h[len(h)-2]
	} else {
		// sequential next
		if st.Index < 0 || st.Index >= len(list) {
			st.Index = 0
		}
		st.Index = (st.Index + 1) % len(list)
		target = list[st.Index]
	}

	if err := b.Set(target); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	st.Last = target
	// update history
	st.History = append(st.History, target)
	if len(st.History) > 50 {
		st.History = st.History[len(st.History)-50:]
	}
	if err := c.Save(st); err != nil {
		fmt.Fprintln(os.Stderr, "warning:", err)
	}
	fmt.Println(target)
}
