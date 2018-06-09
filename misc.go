package main

import (
	"log"
	"os"
	"path"
)

func StrContains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func GetGpmwHome() string {
	xdg_cache_home := os.Getenv("XDG_CACHE_HOME")
	if xdg_cache_home == "" {
		xdg_cache_home = path.Join(os.Getenv("HOME"), ".cache")
	}
	return path.Join(xdg_cache_home, "gpmw")
}

func InitDir(path string) {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalln(err)
		}
	}
}
