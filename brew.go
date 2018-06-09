package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

var (
	GPMW_ETC_DIR   string
	BREW_FILE_PATH string
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

func (p *BrewConfig) MakeBrewFile(_ *cli.Context) error {
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
	file, err := os.Create(BREW_FILE_PATH)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range arr {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return err
}

func UpdateBrew(_ *cli.Context) error {
	err := exec.Command("brew", "upgrade").Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func InstallBrew(_ *cli.Context) error {
	prev, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer os.Chdir(prev)

	os.Chdir(GPMW_ETC_DIR)

	err = exec.Command("brew", "bundle").Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func init() {
	GPMW_ETC_DIR = filepath.Join(GPMW_HOME, "etc")
	InitDir(GPMW_ETC_DIR)
	BREW_FILE_PATH = filepath.Join(GPMW_ETC_DIR, "Brewfile")

	command := cli.Command{
		Name:  "brew",
		Usage: "Make Brewfile, Install or Update Packages Brewfile",
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				cli.ShowCommandHelp(c, "go")
				return nil
			}
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:   "make",
				Usage:  "Make Brewfile",
				Action: CONFIG.Brew.MakeBrewFile,
			},
			{
				Name:   "update",
				Usage:  "brew upgrade",
				Action: UpdateBrew,
			},
			{
				Name:   "install",
				Usage:  "brew install using Brewfile",
				Action: InstallBrew,
			},
		},
	}
	APP.Commands = append(APP.Commands, command)
}
