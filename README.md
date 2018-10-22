[![Build Status](https://travis-ci.com/k-hosokawa/pac.svg?branch=master)](https://travis-ci.com/k-hosokawa/pac)
# pac
Go Package Manager Wrapper for homebrew, go get, and more. Inspired by zplug.

## Install
```sh
$ go get -d github.com/k-hosokawa/pac
```


## Settings
Edit settings in $XDG_CONFIG_HOME/pac/config.toml
```conf.toml
[brew]
tap = ["homebrew/bundle", "homebrew/cask"]
cask = ["mactex", "xquartz"]
pkgs = ["coreutils", "neovim"]
[[brew.optionalPkgs]]
    pkg = "go"
    args = ["--cross-compile-common"]
[go]
repos = [
    "motemen/ghq",
    "peco/peco",
]

[src]
[[src.pkg]]
repo = "zsh-users/zsh-syntax-highlighting"
onZshComp = "src"

[[src.pkg]]
repo = "zsh-users/zsh-syntax-highlighting"
onZshSource = "zsh-syntax-highlighting.zsh"

[[src.pkg]]
repo = "jhawthorn/fzy"
build = ["make"]
onCmd = "fzy"

[[src.pkg]]
repo = "saghul/pythonz"
doClone = false
freeze = true
build = [
    "curl -kLO https://raw.github.com/saghul/pythonz/master/pythonz-install",
    "chmod +x pythonz-install",
    "./pythonz-install",
]
buildEnv = ["PYTHONZ_ROOT=__PACKAGE_HOME__"]
onCmd = "bin/pythonz"
onZshSource = "etc/bashrc"

[[src.pkg]]
repo = "jwilm/alacritty"
onOs = "darwin"
build = [
    "make app",
]
onApp = "target/release/osx/Alacritty.app"

[[src.pkg]]
repo = "b4b4r07/peco-tmux.sh"
onCmd = "peco-tmux.sh"
renameTo = "peco-tmux"
```
### src
* build
    + list of build commands
* buildEnv
    + list of environments when building.
    + you can use `__PACKAGE_HOME__` and `__PAC_HOME__`
    + `__PAC_HOME__`: `$XDG_CACHE_HOME/pac`
    + `__PACKAGE_HOME__`: `$XDG_CACHE_HOME/pac/src/repo/github.com/username/package_dir`
* onCmd
    + create symlink from `__PACKAGE_HOME__/onCmd` to `__PAC_HOME__/bin`
* renameTo
    + rename file when use onCmd
* onApp (for darwin)
    + move to /Application
* onZshComp
    + directory path of zsh completion files
    + add `$XDG_CACHE_HOME/pac/etc/zsh/Completion` to fpath
* onZshSource
    + create symlink from `__PACKAGE_HOME__/onZshSource` to `__PAC_HOME__/etc/zsh/src`
    + use on .zshrc like this
    ```.zshrc
    for f in $XDG_CACHE_HOME/pac/etc/zsh/src/*.zsh; do
        source $f
    done
    ```

## Usage
### example
```sh
$ pac src install
```

### brew
* make
    + make Brewfile
* update
    + brew upgrade
* install
    + brew bundle

### go
* install
    + go get -u repo
* update
    + same as install

### src
* install
* update
