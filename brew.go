package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BrewConfig struct {
	Tap          []string      `toml:tap`
	Cask         []string      `toml:cask`
	Pkgs         []string      `toml:pkgs`
	OptionalPkgs []OptionalPkg `toml:optionalPkgs`
}

type OptionalPkg struct {
	Pkg  string    `toml:pkg`
	Args *[]string `toml:args`
	Link *bool     `toml:link`
}

func (p *OptionalPkg) String() (s string) {
	s = fmt.Sprintf("brew \"%s\"", p.Pkg)
	if p.Args != nil {
		var args []string = *(p.Args)
		for i, v := range args {
			args[i] = fmt.Sprintf("\"%s\"", v)
		}
		s += fmt.Sprintf(", args: [%s]", strings.Join(args, ","))
	}
	if p.Link != nil {
		s += ", link: " + strconv.FormatBool(*(p.Link))
	}
	return
}

func (p *BrewConfig) MakeBrewFile(path string) {
	var arr []string
	if p.Tap != nil {
		for _, v := range p.Tap {
			arr = append(arr, fmt.Sprintf("tap \"%s\"", v))
		}
	}
	if p.Cask != nil {
		for _, v := range p.Cask {
			arr = append(arr, fmt.Sprintf("cask \"%s\"", v))
		}
	}
	if p.Pkgs != nil {
		for _, v := range p.Pkgs {
			arr = append(arr, fmt.Sprintf("brew \"%s\"", v))
		}
	}
	if p.OptionalPkgs != nil {
		for _, v := range p.OptionalPkgs {
			arr = append(arr, v.String())
		}
	}
	// output
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range arr {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
