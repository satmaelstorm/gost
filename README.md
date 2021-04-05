# Go Starter (gost)

The package is designed to quickly start new projects on go.
## Common supported flags
`-v` - verbose output

`-s` - soft launch

`--no-color` - disable colored output
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
However, you can rewrite default aliases and bundles of aliases with your own - gost try to read 
file `gost.aliases.yaml` in current directory to rewrite build-in aliases. Also two flags present:

`--aliases` - file to full rewrite build-in aliases (include loaded from `gost.aliases.yaml`)

`--aliases-add` - file to add and replace build-in aliases

### start
Starts new project in dir `--package-name` in current dir, do `mod init` and all, what do `gost mod` command.
