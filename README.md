[![CircleCI](https://circleci.com/gh/tennashi/got/tree/master.svg?style=shield)](https://circleci.com/gh/tennashi/got/tree/master)

# got
your local machine manager written in Go.  
`got` automates the following tasks:
```bash
$ git clone https://github.com/your_name/dotfiles
$ ln -s /path/to/your/dotfiles/bashrc ~/.bashrc
$ ln -s /path/to/your/dotfiles/xmonad/ ~/.xmonad
...
$ sudo apt install neovim xmonad fish ...
```

# Usage
## Install
Get from the [release page](https://github.com/tennashi/got/releases)  
or
```bash
$ go get -u http://github.com/tennashi/got
```

## Sync your dotfiles
```bash
$ got sync http://github.com/your_name/dotfiles /path/to/your/dotfiles

# or if you want to overwrite config.(default: ~/.config/got/config.toml)
# It's recomended when using `got` for the first time.
$ got sync -w http://github.com/your_name/dotfiles /path/to/our/dotfiles

# If you wrote `config.toml`, you can omit the dotfiles config.
$ got sync
```
This command execute following tasks.
  * `git clone http://github.com/your_name/dotfiles /path/to/your/dotfiles`
  * Install packages specified by `Gotfile.toml`
  * Make symlinks specified by `Gotfile.toml`
    * If a file exist in destination of symlink, skip creating the symlink.

## Add package
```bash
$ got get apt sl
```
This command installs `sl` as `apt` command and append the `sl` config to `Gotfile.toml`.
```toml
...
[package]
  [package.apt]
    names = [
      "sl", # Autometional added
    ]
```

If you set `default_manager` in `Gotfile.toml`, you can omit the manager you use.
```bash
$ got get sl
```
and `config.toml`
```toml
default_manager = "apt"
```
This situation is the same as running `got get apt sl`.


## Config
Set the location of your dotfiles repository in the config file.

Load the config file in the following order.
  * config file specifyed by `-c` option
  * $XDG_CONFIG_HOME/got/config.toml
  * xdg_config_dir/got/config.toml (where xdg_config_dir is listed in $XDG_CONFIG_DIRS)
  * ~/.config/got/config.toml

### Example
```toml
[dotfiles]
  local = "/path/to/your/dotfiles" # your local dotfiles directory
  remote = "https://github.com/your_name/dotfiles" # your remote dotfiles repository

default_manager = "apt" # default package manager
```

## Gotfile
Set source and destination of symlink in `Gotfile.toml`.  
Put gotfile in `/path/to/your/dotfiles/Gotfile.toml`.

### Example
```toml
[[dotfile]]
  dest = "~/.bashrc" # Set symlink destination.
  src = "bashrc"     # Set symlink source.
  
[[dotfile]]
  dest = "~/.xmonad" # It distinguishes between directories and files automatically.
  src = "xmonad"

[package]
  [package.apt]
    names = [
      "sl",
      "neovim",
      "fish",
    ]
```
