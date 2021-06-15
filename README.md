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
