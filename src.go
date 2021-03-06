package main

import (
	"fmt"
	"log"
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
	Repo        string    `toml:repo`
	DoClone     *bool     `toml:doClone`
	OnOs        *string   `toml:onOs`
	Freeze      *bool     `toml:freeze`
	Build       *[]string `toml:build`
	BuildEnv    *[]string `toml:buildEnv`
	OnApp       *string   `toml:onApp`
	OnCmd       *string   `toml:onCmd`
	OnZshSource *string   `toml:onZshSource`
	OnZshComp   *string   `toml:onZshComp`
	RenameTo    *string   `toml:renameTo`
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
	r := strings.NewReplacer(
		"__PAC_HOME__", PAC_HOME,
		"__PACKAGE_HOME__", package_path,
	)
	for _, env := range *env_arr {
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

func install_bin(dir_path string, cmd *string, rename *string) {
	bin_path, err := filepath.Abs(filepath.Join(dir_path, *cmd))
	if err != nil {
		fmt.Println(err)
	}
	if err := os.Chmod(bin_path, 0755); err != nil {
		fmt.Println(err)
	}
	cmd_path := *cmd
	if rename != nil {
		cmd_path = *rename
	}
	sympath := filepath.Join(PAC_HOME, "bin", filepath.Base(cmd_path))
	if err = RmIfExist(sympath); err != nil {
		fmt.Println(nil)
	}

	if err := os.Symlink(bin_path, sympath); err != nil {
		fmt.Println(err)
	}
}

func install_zsh_completion(repo string, comp_path string) {
	dist := filepath.Join(ZSH_COMP_DIR, filepath.Base(repo))
	comp_path = filepath.Join(GetRepoPath(repo), comp_path)
	if _, err := os.Stat(dist); !os.IsNotExist(err) {
		return
	}
	if f, err := os.Stat(comp_path); err != nil {
		fmt.Println(err)
	} else {
		if !f.IsDir() {
			log.Fatalln("onZshComp must be directory")
			return
		}
	}

	if err := os.Symlink(comp_path, dist); err != nil {
		fmt.Println(err)
	}
}

func install_zsh_source(repo string, source_path string) {
	dist := filepath.Join(ZSH_SOURCE_DIR, filepath.Base(repo)+".zsh")
	if _, err := os.Stat(dist); !os.IsNotExist(err) {
		return
	}

	source_path, err := filepath.Abs(
		filepath.Join(GetRepoPath(repo), source_path),
	)
	if err != nil {
		fmt.Println(nil)
	}

	if err = os.Symlink(source_path, dist); err != nil {
		fmt.Println(err)
	}
}

func install_src(pkg SrcPkg, isFin chan bool, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(pkg.Repo)
	if pkg.OnOs != nil {
		if *(pkg.OnOs) != runtime.GOOS {
			isFin <- true
			return
		}
	}
	dir_path := filepath.Join(PAC_HOME, "src", "github.com", pkg.Repo)

	if pkg.DoClone == nil || *(pkg.DoClone) {
		clone(dir_path, pkg.Repo)
	} else {
		InitDir(dir_path)
	}

	if pkg.Build != nil {
		build(dir_path, pkg.Build, pkg.BuildEnv)
	}

	if pkg.OnCmd != nil {
		install_bin(dir_path, pkg.OnCmd, pkg.RenameTo)
	}

	if pkg.OnZshComp != nil {
		install_zsh_completion(pkg.Repo, *(pkg.OnZshComp))
	}

	if pkg.OnZshSource != nil {
		install_zsh_source(pkg.Repo, *(pkg.OnZshSource))
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

func check_update(repo_path string) bool {
	prev, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Chdir(prev)

	os.Chdir(repo_path)
	out, err := exec.Command("git", "pull").Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if string(out) == "Already up to date.\n" {
		return false
	} else {
		return true
	}
}

func update_src(src SrcPkg, isFin chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	_os := src.OnOs != nil && *(src.OnOs) != runtime.GOOS
	_freeze := src.Freeze != nil && *(src.Freeze)
	_clone := src.DoClone != nil && !*(src.DoClone)
	if _os || _freeze || _clone {
		isFin <- true
		return
	}
	repo_path := GetRepoPath(src.Repo)
	if check_update(repo_path) {
		if src.Build != nil {
			build(repo_path, src.Build, src.BuildEnv)
		}
		fmt.Printf("%-50s %s\n", src.Repo, "Update !!")
	} else {
		fmt.Printf("%-50s %s\n", src.Repo, "not update")
	}
	isFin <- true
}

func UpdateSrc(_ *cli.Context) error {
	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(CONFIG.Src.Pkg))
	for _, src := range CONFIG.Src.Pkg {
		wg.Add(1)
		update_src(src, isFin, wg)
	}
	wg.Wait()
	close(isFin)
	fmt.Println("Updated Packages from Source")
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
