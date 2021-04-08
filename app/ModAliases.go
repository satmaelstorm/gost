package app

import (
	_ "embed" //golint: import embed
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	packageTypeAlias = iota
	packageTypeGitHub
	packageTypeOther
)

const defaultFileName = "gost.aliases.yaml"

//go:embed gost.aliases.yaml
var defAliases []byte
var defModAliases ModAliases

//PackageName - type for Package Names
type PackageName string

func (p PackageName) getPackageType() int {
	s := strings.Split(string(p), "/")
	switch len(s) {
	case 0:
		return packageTypeAlias
	case 1:
		return packageTypeAlias
	case 2:
		f := strings.ToLower(s[0])
		if "gopkg.in" == f || "golang.org" == f {
			return packageTypeOther
		}
		return packageTypeGitHub
	default:
		return packageTypeOther
	}
}

//ModAliases - represents aliases and bundles set
type ModAliases struct {
	Packages map[string]string        `yaml:"packages"`
	Bundles  map[string][]PackageName `yaml:"bundles"`
}

func init() {
	f, err := ioutil.ReadFile(defaultFileName)
	if err == nil {
		fmt.Println("Load aliases from " + defaultFileName)
		defAliases = f
	}
	ma, err := getAliasesByBytes(defAliases)
	if err != nil {
		fmt.Println("[ERROR]: " + err.Error())
		os.Exit(1)
	}
	defModAliases = ma
	err = defModAliases.validate()
	if err != nil {
		fmt.Println("[ERROR]: " + err.Error())
		os.Exit(1)
	}
}

func getAliasesByBytes(str []byte) (ModAliases, error) {
	var result ModAliases
	err := yaml.Unmarshal(str, &result)
	return result, err
}

//GetAliases - return current default aliases and bundles set
func GetAliases() ModAliases {
	return defModAliases
}

//GetDefaultAliasesH - return string of default aliases and bundles set for description of commands
func GetDefaultAliasesHelp() string {
	return string(defAliases)
}

func getAliasesByFile(fileName string, isValidate bool) (ModAliases, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ModAliases{}, err
	}
	ma, err := getAliasesByBytes(f)
	if err != nil {
		return ModAliases{}, err
	}
	if isValidate {
		err := ma.validate()
		if err != nil {
			return ModAliases{}, err
		}
	}
	return ma, nil
}

//LoadAliasesFromFlags - load aliases and bundles set from parameters from command line
func LoadAliasesFromFlags(outIo, errIo io.Writer, aliasesFile, addAliasesFile string) bool {
	if aliasesFile != "" {
		ma, err := getAliasesByFile(aliasesFile, true)
		if err != nil {
			writeError(errIo, err.Error())
			return false
		}
		defModAliases = ma
		_, _ = outIo.Write([]byte("Load aliases from " + aliasesFile))
	}
	if addAliasesFile != "" {
		ma, err := getAliasesByFile(addAliasesFile, false)
		if err != nil {
			writeError(errIo, err.Error())
			return false
		}
		newMa, err := defModAliases.Glue(ma)
		if err != nil {
			writeError(errIo, err.Error())
			return false
		}
		defModAliases = newMa
		_, _ = outIo.Write([]byte("Load additional aliases from " + addAliasesFile))
	}
	return true
}

func (m ModAliases) validate() error {
	for name, packages := range m.Bundles {
		if _, ok := m.Packages[name]; ok {
			return errors.New("has package with name " + name + ", bundle with some name is not allowed")
		}
		for _, pkg := range packages {
			if pkg.getPackageType() != packageTypeAlias {
				continue
			}
			_, ok1 := m.Packages[string(pkg)]
			_, ok2 := m.Bundles[string(pkg)]
			if !ok1 && !ok2 {
				return errors.New("bundle " + name + "require package or bundle alias " +
					string(pkg) + " but it not present in list of packages or bundles")
			}
		}
	}
	return nil
}

//Glue - glue two ModAliases in one
func (m ModAliases) Glue(src ModAliases) (ModAliases, error) {
	for name, pkg := range src.Packages {
		m.Packages[name] = pkg
	}
	for name, pkg := range src.Bundles {
		m.Bundles[name] = pkg
	}
	err := m.validate()
	if err != nil {
		return ModAliases{}, err
	}
	return m, nil
}
