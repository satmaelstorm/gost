Tired of remembering the names of the libraries you need on github? Tired of writing long commands like `go get -u github.com/gomodule/redigo/redis`? 
Then the `gost` is for you!


# Go Starter (gost)
[![Go Report](https://goreportcard.com/badge/github.com/satmaelstorm/gost)](https://goreportcard.com/report/github.com/satmaelstorm/gost)
[![GoDoc](https://godoc.org/github.com/satmaelstorm/gost?status.svg)](http://godoc.org/github.com/satmaelstorm/gost)
[![Coverage Status](https://coveralls.io/repos/github/satmaelstorm/gost/badge.svg?branch=master)](https://coveralls.io/github/satmaelstorm/gost?branch=master)
![Go](https://github.com/satmaelstorm/gost/workflows/Go/badge.svg)

The command allows you to assign aliases for popular libraries and not remember long repository names.
## Install
Require go 1.16.0 or later, to install run:

`go install github.com/satmaelstorm/gost@latest`
## Common supported flags
`-v` - verbose output

`-s` - soft launch

`--no-color` - disable colored output

## Rewrite Aliases
You can export env `GOST_ALIASES` to rewrite build-in aliases 
and export `GOST_ADD_ALIASES` to add and replace build-in aliases. For example: `GOST_ALIASES="/home/user/gost.aliases.yaml" ./gost mod -s webserver`

You can rewrite default aliases and bundles of aliases with your own - gost try to read 
file `gost.aliases.yaml` in a current directory to rewrite build-in aliases. 

Also two flags present:

`--aliases` - file to full rewrite build-in aliases (include loaded from `gost.aliases.yaml`)

`--aliases-add` - file to add and replace build-in aliases

Reading and replacing order:
1. From env
2. From file in current directory
3. From flags
 
## Multi-threads (experimental)
You can run `ghost mod` with` --threads = N` for multithreading, but then the order in which 
the `go get` is executed is not guaranteed, there may be side effects.

## Supported commands
### mod
Do `go get -u` commands with use of aliases and bundles of aliases.

Usage: `gost mod module1 module2 ...`

Sample (with soft launch):
```bash
gost mod webserver -s
    
Use soft Launch

/usr/local/go/bin/go get -u github.com/valyala/fasthttp
/usr/local/go/bin/go get -u github.com/fasthttp/router
/usr/local/go/bin/go get -u github.com/stretchr/testify
/usr/local/go/bin/go get -u github.com/satmaelstorm/envviper
/usr/local/go/bin/go mod download
```
If there is strictly one `/` in the name of the module - it is perceived as a github module:
```bash
gost mod tidwall/rtree
Use soft Launch

/usr/local/go/bin/go get -u github.com/tidwall/rtree
/usr/local/go/bin/go mod download
```
But gost recognize `gopkg.in` and `golang.org`:
```bash
gost mod gopkg.in/yaml.v1 -s
Use soft Launch

/usr/local/go/bin/go get -u gopkg.in/yaml.v1
/usr/local/go/bin/go mod download
```
Names with more than one `/` - are sent to go get unchanged:
```bash
gost mod gopkg.in/Graylog2/go-gelf.v1 -s
Use soft Launch

/usr/local/go/bin/go get -u gopkg.in/Graylog2/go-gelf.v1
/usr/local/go/bin/go mod download
```

### start
Starts new project in dir `--package-name` in current dir, do `mod init` and all, what do `gost mod` command.
