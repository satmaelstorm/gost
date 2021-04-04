package app

import (
	_ "embed"
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

const (
	packageTypeAlias = iota
	packageTypeGitHub
	packageTypeOther
)

//go:embed gost.aliases.yaml
var defAliases []byte
var defModAliases ModAliases

type ExecutorFunc func(string, int) error
type PackageName string

func (p PackageName) getPackageType() int {
	switch len(strings.Split(string(p), "/")) {
	case 0:
		return packageTypeAlias
	case 1:
		return packageTypeAlias
	case 2:
		return packageTypeGitHub
	default:
		return packageTypeOther
	}
}

type ModAliases struct {
	Packages map[string]string        `yaml:"packages"`
	Bundles  map[string][]PackageName `yaml:"bundles"`
}

func init() {
	err := yaml.Unmarshal(defAliases, &defModAliases)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = defModAliases.validate()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func GetDefaultAliases() ModAliases {
	return defModAliases
}

func GetDefaultAliasesHelp() string {
	return string(defAliases)
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
