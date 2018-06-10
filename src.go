package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

type SrcConfig struct {
	Pkg []SrcPkg `toml:pkg`
}

type SrcPkg struct {
	Repo     string    `toml:repo`
	DoClone  *bool     `toml:doClone`
	OnOs     *string   `toml:onOs`
	Freeze   *bool     `toml:freeze`
	Build    *[]string `toml:build`
	BuildEnv *[]string `toml:buildEnv`
	OnApp    *string   `toml:onApp`
	OnCmd    *string   `toml:onCmd`
}

func clone(dir_path string, repo string) {
	if f, err := os.Stat(dir_path); os.IsNotExist(err) || !f.IsDir() {
		err := exec.Command(
			"git", "clone",
			"https://github.com/"+repo,
			dir_path,
		).Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func convert_env(env_arr *[]string, package_path string) (new_arr []string) {
	for _, env := range *env_arr {
		r := strings.NewReplacer(
			"__GPMW_HOME__", GPMW_HOME,
			"__PACKAGE_HOME__", package_path,
		)
		new_arr = append(new_arr, r.Replace(env))
	}
	return
}

func build(dir_path string, build_arr *[]string, env_arr *[]string) {
	prev, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Chdir(prev)

	os.Chdir(dir_path)
	for _, cmd_str := range *build_arr {
		cmds := strings.Split(cmd_str, " ")
		cmd := exec.Command(cmds[0], cmds[1:]...)
		if env_arr != nil {
			new_arr := convert_env(env_arr, dir_path)
			cmd.Env = append(os.Environ(), new_arr...)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func install_bin(dir_path string, cmd *string) {
	bin_path, err := filepath.Abs(filepath.Join(dir_path, *cmd))
	if err != nil {
		fmt.Println(err)
	}
	if err := os.Chmod(bin_path, 0755); err != nil {
		fmt.Println(err)
	}
	sympath := filepath.Join(GPMW_HOME, "bin", filepath.Base(*cmd))
	if _, err := os.Stat(sympath); !os.IsNotExist(err) {
		err = os.Remove(sympath)
		if err != nil {
			fmt.Println(nil)
		}
	}

	if err := os.Symlink(bin_path, sympath); err != nil {
		fmt.Println(err)
	}
}

func install_src(
	src SrcPkg,
	isFin chan bool,
	wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(src.Repo)
	if src.OnOs != nil {
		if *(src.OnOs) != runtime.GOOS {
			return
		}
	}
	dir_path := filepath.Join(GPMW_HOME, "src", "github.com", src.Repo)

	if src.DoClone == nil || *(src.DoClone) {
		clone(dir_path, src.Repo)
	} else {
		InitDir(dir_path)
	}

	if src.Build != nil {
		build(dir_path, src.Build, src.BuildEnv)
	}

	if src.OnCmd != nil {
		install_bin(dir_path, src.OnCmd)
	}
	isFin <- true
}

func InstallSrc(_ *cli.Context) error {
	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(CONFIG.Src.Pkg))
	for _, src := range CONFIG.Src.Pkg {
		wg.Add(1)
		install_src(src, isFin, wg)
	}
	wg.Wait()
	close(isFin)
	fmt.Println("Installed Packages from Source")
	return nil
}

// TODO
func UpdateSrc(_ *cli.Context) error {
	return nil
}

func init() {
	command := cli.Command{
		Name:  "src",
		Usage: "Install or Update Packages from source",
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				cli.ShowCommandHelp(c, "src")
				return nil
			}
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:   "update",
				Usage:  "Update Installed Packages",
				Action: UpdateSrc,
			},
			{
				Name:   "install",
				Usage:  "Install Packages from source",
				Action: InstallSrc,
			},
		},
	}
	APP.Commands = append(APP.Commands, command)
}
