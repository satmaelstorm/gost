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

func (m ModAliases) validate() error {
	for name, packages := range m.Bundles {
		if _, ok := m.Packages[name]; ok {
			return errors.New("has package with name " + name + ", bundle with some name is not allowed")
		}
		for _, pkg := range packages {
			if pkg.getPackageType() != packageTypeAlias {
				continue
			}
			if _, ok := m.Packages[string(pkg)]; !ok {
				return errors.New("bundle " + name + "require package alias " +
					string(pkg) + " but it not present in list of packages")
			}
		}
	}
	return nil
}
