package app

import (
	_ "embed"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
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

type GoGetModules struct {
	visitedAliases map[string]bool
	goGet          ExecutorFunc
	goMod          ExecutorFunc
	verboseLevel   int
}

func (g *GoGetModules) AsSoftLaunch() *GoGetModules {
	g.goGet = goGetTest
	g.goMod = goModSoft
	return g
}

func (g *GoGetModules) AsExec() *GoGetModules {
	g.goGet = goGetSimple
	g.goMod = goModExec
	return g
}

func (g *GoGetModules) VerboseLevel(v int) *GoGetModules {
	g.verboseLevel = v
	return g
}

func (g *GoGetModules) Run(names []string, aliases ModAliases) {
	if nil == g.goGet {
		g.goGet = goGetSimple
		g.goMod = goModExec
	}
	g.run(names, aliases)
	err := g.goMod("download", g.verboseLevel)
	if err != nil {
		fmt.Println(err)
	}
}

func (g *GoGetModules) run(names []string, aliases ModAliases) {
	if nil == g.goGet {
		log.Fatal("undefined goGet function")
	}
	g.visitedAliases = make(map[string]bool)
	for _, n := range names {
		err := g.goGetFull(n, aliases)
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *GoGetModules) getNameAndVersion(in string) (string, string) {
	s := strings.Split(in, "@")
	if len(s) < 1 {
		return "", ""
	}
	if len(s) == 1 {
		return s[0], ""
	}
	return s[0], s[1]
}

func (g *GoGetModules) nameWithVersion(name, ver string) string {
	if ver != "" {
		return name + "@" + ver
	}
	return name
}

func (g *GoGetModules) goGetFull(name string, aliases ModAliases) error {
	n, v := g.getNameAndVersion(name)
	switch PackageName(n).getPackageType() {
	case packageTypeAlias:
		return g.goGetAliases(n, v, aliases)
	case packageTypeGitHub:
		return g.goGetGitHub(n, v)
	case packageTypeOther:
		return g.goGetOther(n, v)
	}
	return errors.New("unknown type of package name " + name)
}

func (g *GoGetModules) goGetAliases(name, ver string, aliases ModAliases) error {
	if _, ok := g.visitedAliases[name]; ok {
		return nil
	}
	g.visitedAliases[name] = true
	if bundle, ok := aliases.Bundles[name]; ok {
		if len(bundle) > 0 {
			for _, pkg := range bundle {
				err := g.goGetAliases(string(pkg), "", aliases)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	if pkg, ok := aliases.Packages[name]; ok {
		return g.goGet(g.nameWithVersion(pkg, ver), g.verboseLevel)
	}
	return errors.New("not found package or bundle " + name)
}

func (g *GoGetModules) goGetGitHub(name string, ver string) error {
	return g.goGet("github.com/" + g.nameWithVersion(name, ver), g.verboseLevel)
}

func (g *GoGetModules) goGetOther(name string, ver string) error {
	return g.goGet(g.nameWithVersion(name, ver), g.verboseLevel)
}

func goGetSimple(name string, verboseLevel int) error {
	c := exec.Command("go", "get", "-u", name)
	if verboseLevel > 0 {
		log.Println("Execute: " + c.String())
	}
	out, err := c.CombinedOutput()
	if verboseLevel > 0 {
		log.Println(string(out))
	}
	return err
}

func goGetTest(name string, verboseLevel int) error {
	c := exec.Command("go", "get", "-u", name)
	fmt.Println(c)
	return nil
}

func goMod(sub string) *exec.Cmd{
	return exec.Command("go", "mod", sub)
}

func goModExec(sub string, verboseLevel int) error {
	c := goMod(sub)
	if verboseLevel > 0 {
		log.Println("Executing: " + c.String())
	}
	out, err := c.CombinedOutput()
	if verboseLevel > 0 {
		log.Println(string(out))
	}
	return err
}

func goModSoft(sub string, verboseLevel int) error {
	c := goMod(sub)
	fmt.Println(c.String())
	return nil
}
