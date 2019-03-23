# got
dotfiles manager written in Go.  
`got` automates the following tasks:
```bash
$ git clone https://github.com/your_name/dotfiles
$ ln -s /path/to/your/dotfiles/bashrc ~/.bashrc
$ ln -s /path/to/your/dotfiles/xmonad/ ~/.xmonad
...
```

# Usage
## Install
```bash
$ go get -u http://github.com/tennashi/got
```

## Sync your dotfiles
```bash
$ got sync http://github.com/your_name/dotfiles /path/to/your/dotfiles

# or if you want to overwrite config(default: ~/.config/got/config.toml)
$ got sync -w http://github.com/your_name/dotfiles /path/to/our/dotfiles
```
If a file exist in destination of symlink, skip creating the symlink.

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
  local = "/path/to/your/dotfiles"
  remote = "https://github.com/your_name/dotfiles"

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
```
