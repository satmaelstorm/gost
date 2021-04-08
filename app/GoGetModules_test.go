package app

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"testing"
)

const testAliasesConfig = `
packages:
  fasthttp: github.com/valyala/fasthttp
  fastrouter: github.com/fasthttp/router
bundles:
  webserver:
    - fasthttp
    - fastrouter
    - satmaelstorm/bitmap@v1.2.0
    - bitbucket.org/satmaelstorm/bitmap
`

type goGetModulesTestSuite struct {
	suite.Suite
	g *GoGetModules
	aliases ModAliases
}

func TestGoGetModules(t *testing.T) {
	suite.Run(t, new(goGetModulesTestSuite))
}

func (s *goGetModulesTestSuite) SetupSuite() {
	a, err := getAliasesByBytes([]byte(testAliasesConfig))
	if err != nil {
		s.T().Fatal(err.Error())
	}
	s.aliases = a
}

func (s *goGetModulesTestSuite) SetupTest() {
	s.g = new(GoGetModules)
	s.g.stdOut = new(bytes.Buffer)
	s.g.errorOut = new(bytes.Buffer)
	s.g.visitedAliases = make(map[string]bool)
}

func (s *goGetModulesTestSuite) TearDownTest() {
	s.g = nil
}

func (s *goGetModulesTestSuite) TestGoGetGitHub() {
	s.g.goGetGitHub("satmaelstorm/bitmap", "v1.2.0")
	s.Equal(1, len(s.g.commands))
	s.NotNil(s.g.commands[0])
	cmd := s.g.commands[0].String()
	s.Contains(cmd, "go get -u github.com/satmaelstorm/bitmap@v1.2.0")
}

func (s *goGetModulesTestSuite) TestGoGetOther() {
	s.g.goGetOther("bitbucket.org/satmaelstorm/bitmap", "v1.2.0")
	s.Equal(1, len(s.g.commands))
	s.NotNil(s.g.commands[0])
	cmd := s.g.commands[0].String()
	s.Contains(cmd, "go get -u bitbucket.org/satmaelstorm/bitmap@v1.2.0")
}

func (s *goGetModulesTestSuite) TestNameWithVersion() {
	s.Equal("satmaelstorm/bitmap", s.g.nameWithVersion("satmaelstorm/bitmap", ""))
	s.Equal("satmaelstorm/bitmap@v1.2.0", s.g.nameWithVersion("satmaelstorm/bitmap", "v1.2.0"))
}

func (s *goGetModulesTestSuite) TestGetNameAndVersion() {
	var r1,r2 string
	r1,r2 = s.g.getNameAndVersion("")
	s.Equal(r1, "")
	s.Equal(r2, "")
	r1,r2 = s.g.getNameAndVersion("satmaelstorm/bitmap")
	s.Equal(r1, "satmaelstorm/bitmap")
	s.Equal(r2, "")
	r1,r2 = s.g.getNameAndVersion("satmaelstorm/bitmap@v1.2.0")
	s.Equal(r1, "satmaelstorm/bitmap")
	s.Equal(r2, "v1.2.0")
}

func (s *goGetModulesTestSuite) TestAsSoftLaunch() {
	s.False(s.g.isSoft)
	s.g.AsSoftLaunch()
	s.True(s.g.isSoft)
}

func (s *goGetModulesTestSuite) TestVerboseLevel() {
	s.Equal(0, s.g.verboseLevel)
	s.g.VerboseLevel(1)
	s.Equal(1, s.g.verboseLevel)
}

func (s *goGetModulesTestSuite) TestGoGetAliases() {
	var err error
	err = s.g.goGetAliases("bitmap","", s.aliases)
	s.NotNil(err)
	err = s.g.goGetAliases("bitmap","", s.aliases)
	s.Nil(err)
	s.Equal(0, len(s.g.commands))
	err = s.g.goGetAliases("fasthttp", "v1.20.0", s.aliases)
	s.Nil(err)
	s.Equal(1, len(s.g.commands))
	s.NotNil(s.g.commands[0])
	s.Contains(s.g.commands[0].String(), "go get -u github.com/valyala/fasthttp@v1.20.0")
	err = s.g.goGetAliases("webserver", "", s.aliases)
	s.Nil(err)
	s.Equal(4, len(s.g.commands))
	s.NotNil(s.g.commands[1])
	s.NotNil(s.g.commands[2])
	s.NotNil(s.g.commands[3])
	s.Contains(s.g.commands[1].String(), "go get -u github.com/fasthttp/router")
	s.Contains(s.g.commands[2].String(), "go get -u github.com/satmaelstorm/bitmap@v1.2.0")
	s.Contains(s.g.commands[3].String(), "go get -u bitbucket.org/satmaelstorm/bitmap")
}

func (s *goGetModulesTestSuite) TestGoGetFull() {
	err := s.g.goGetAliases("webserver", "", s.aliases)
	s.Nil(err)
	s.Equal(4, len(s.g.commands))
	s.NotNil(s.g.commands[0])
	s.NotNil(s.g.commands[1])
	s.NotNil(s.g.commands[2])
	s.NotNil(s.g.commands[3])
	s.Contains(s.g.commands[0].String(), "go get -u github.com/valyala/fasthttp")
	s.Contains(s.g.commands[1].String(), "go get -u github.com/fasthttp/router")
	s.Contains(s.g.commands[2].String(), "go get -u github.com/satmaelstorm/bitmap@v1.2.0")
	s.Contains(s.g.commands[3].String(), "go get -u bitbucket.org/satmaelstorm/bitmap")
}

func (s *goGetModulesTestSuite) TestRunWithError() {
	s.g.run([]string{"cli"}, s.aliases)
	eo := s.g.errorOut.(*bytes.Buffer)
	s.True(len(eo.String()) > 0)
}

func (s *goGetModulesTestSuite) TestRun() {
	s.g.run([]string{"webserver"}, s.aliases)
	eo := s.g.errorOut.(*bytes.Buffer)
	s.Equal(0, len(eo.String()))
	s.Equal(4, len(s.g.commands))
	s.NotNil(s.g.commands[0])
	s.NotNil(s.g.commands[1])
	s.NotNil(s.g.commands[2])
	s.NotNil(s.g.commands[3])
	s.Contains(s.g.commands[0].String(), "go get -u github.com/valyala/fasthttp")
	s.Contains(s.g.commands[1].String(), "go get -u github.com/fasthttp/router")
	s.Contains(s.g.commands[2].String(), "go get -u github.com/satmaelstorm/bitmap@v1.2.0")
	s.Contains(s.g.commands[3].String(), "go get -u bitbucket.org/satmaelstorm/bitmap")
}


