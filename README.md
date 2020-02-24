# got
Tool for managing commands written in Go.
This tool was inspired by [Manage Go tools via Go modules](https://marcofranssen.nl/manage-go-tools-via-go-modules/)

# Install
```bash
$ go get http://github.com/tennashi/got/cmd/got
```

# Usage
```
got - go packages manager

Usage:
  got command [arguments]

Commands:
  help
    print this help

  version
    print got command version

  get [-lu] [-c command] package
    install the package

  remove [package|command]
    remove the package

```

## Install command written in Go
```
$ got get github.com/tennashi/got # == go get github.com/tennashi/got
$ got get tennashi/got # == go get github.com/tennashi/got
$ got get -c got tennashi/got # == go get github.com/tennashi/got/cmd/got
$ got get -c hoge tennashi/got # == go get github.com/tennashi/got/cmd/hoge
```

## Update all installed command
```
$ got get -l -u
```

## Remove command
```
$ got remove github.com/tennashi/got/cmd/got
or
$ got remove got
```
These commands remove `$GOBIN/got` (If unset `$GOBIN`, it defaults `$GOPATH/bin`.) and remove the import `github.com/tennashi/got/cmd/got` from `~/.local/share/got/tools.go`.
