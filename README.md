# got
Tool for managing commands written in Go.

## Require
* `go >= 1.16`
  * using `go install` command to install a specific version

## Install
```bash
$ go get http://github.com/tennashi/got/cmd/got
```

## Usage
```
NAME:
   got - package manager for commands written in Go

USAGE:
   got [global options] command [command options] [arguments...]

COMMANDS:
   install  install the specified package
   upgrade  upgrade installed packages
   list     list installed packages
   enable   enable the specified executable
   disable  disable the specified executable
   show     show the specified package
   pin      pin the version of specified package
   unpin    unpin the version of specified package
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  config file [$GOT_CONFIG_FILE]
   --debug                   debug mode (default: false) [$GOT_DEBUG]
   --help, -h                show help (default: false)
```

### Install the command written in Go
```
$ got install github.com/tennashi/got # == go install github.com/tennashi/got/...@latest
$ got install tennashi/got # == go install github.com/tennashi/got/...@latest
```

### Upgrade installed commands
```
$ got upgrade
```

### List installed commands
```
$ got list
```

## Configuration
Priority is given in the following order.

* `--config`/`-c` command line option
* `$GOT_CONFIG_FILE`
* `$XDG_CONFIG_HOME/got/config.toml`
  * default(Linux): `$HOME/.config/got/config.toml`

```toml
install_all = true
bin_dir = "$GOBIN" # if $GOBIN is not defined, "$GOPATH/bin" is used
data_dir = "$XDG_DATA_HOME/got" # default: "$HOME/.local/share/got
```

* `install_all`: Specifies whether to install all main packages under the specified module path
* `bin_dir`: Directory to create symlinks for executables
* `data_dir`: Directory of data to manage the installed modules and the executable files
