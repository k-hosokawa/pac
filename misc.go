package main

import (
	"log"
	"os"
	"path/filepath"
)

func GetGpmwHome() string {
	xdg_cache_home := os.Getenv("XDG_CACHE_HOME")
	if xdg_cache_home == "" {
		xdg_cache_home = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return filepath.Join(xdg_cache_home, "gpmw")
}

func InitDir(path string) {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalln(err)
		}
	}
}
