package app

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

//GoGetModules - executor for gost mod command
type GoGetModules struct {
	visitedAliases map[string]bool
	isSoft         bool
	verboseLevel   int
	commands       []*exec.Cmd
	errorOut       io.Writer
	stdOut         io.Writer
}

//AsSoftLaunch - prepare executor to soft launch
func (g *GoGetModules) AsSoftLaunch() *GoGetModules {
	g.isSoft = true
	return g
}

//VerboseLevel - set verbose level of output (0 or 1)
func (g *GoGetModules) VerboseLevel(v int) *GoGetModules {
	g.verboseLevel = v
	return g
}

//Run - runs executor
func (g *GoGetModules) Run(names []string, aliases ModAliases) {
	g.errorOut = os.Stderr
	g.stdOut = os.Stdout
	g.run(names, aliases)
	modCmd := goMod("download")
	g.commands = append(g.commands, modCmd)
	execCommands(g.stdOut, g.errorOut, g.commands, g.verboseLevel, g.isSoft)
}

func (g *GoGetModules) run(names []string, aliases ModAliases) {
	g.visitedAliases = make(map[string]bool)
	for _, n := range names {
		err := g.goGetFull(n, aliases)
		if err != nil {
			_, _ = g.errorOut.Write([]byte(err.Error() + "\n"))
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
		g.goGetGitHub(n, v)
		return nil
	case packageTypeOther:
		g.goGetOther(n, v)
		return nil
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
				err := g.goGetFull(string(pkg), aliases)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	if pkg, ok := aliases.Packages[name]; ok {
		g.commands = append(g.commands, goGetWithUpdate(g.nameWithVersion(pkg, ver)))
		return nil
	}
	return errors.New("not found package or bundle " + name)
}

func (g *GoGetModules) goGetGitHub(name string, ver string) {
	g.commands = append(g.commands, goGetWithUpdate("github.com/"+g.nameWithVersion(name, ver)))
}

func (g *GoGetModules) goGetOther(name string, ver string) {
	g.commands = append(g.commands, goGetWithUpdate(g.nameWithVersion(name, ver)))
}
