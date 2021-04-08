package app

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type packageNameTestSuite struct {
	suite.Suite
}

func TestPackageName(t *testing.T) {
	suite.Run(t, new(packageNameTestSuite))
}

func (s *packageNameTestSuite) TestGetPackageType() {
	s.Equal(packageTypeAlias, PackageName("webserver").getPackageType())
	s.Equal(packageTypeGitHub, PackageName("satmaelstorm/bitmap").getPackageType())
	s.Equal(packageTypeOther, PackageName("gopkg.in/yaml.v1").getPackageType())
	s.Equal(packageTypeOther, PackageName("golang.org/sync").getPackageType())
	s.Equal(packageTypeOther, PackageName("golang.org/x/sync").getPackageType())
	s.Equal(packageTypeOther, PackageName("bitbucket.org/satmaelstorm/bitmap").getPackageType())
}

type modAliasesTestSuite struct {
	suite.Suite
}

func TestModAliases(t *testing.T) {
	suite.Run(t, new(modAliasesTestSuite))
}

func (s *modAliasesTestSuite) TestGetByBytes() {
	a, err := getAliasesByBytes([]byte(testAliasesConfig))
	s.Nil(err)
	s.Equal(2, len(a.Packages))
	s.Equal(1, len(a.Bundles))
	s.Contains(a.Packages, "fasthttp")
	s.Contains(a.Packages, "fastrouter")
	s.Contains(a.Bundles, "webserver")
	s.Equal(4, len(a.Bundles["webserver"]))
	s.Equal("fasthttp", string(a.Bundles["webserver"][0]))
	s.Equal("fastrouter", string(a.Bundles["webserver"][1]))
	s.Equal("satmaelstorm/bitmap@v1.2.0", string(a.Bundles["webserver"][2]))
	s.Equal("bitbucket.org/satmaelstorm/bitmap", string(a.Bundles["webserver"][3]))
}

func (s *modAliasesTestSuite) TestValidate() {
	a, err := getAliasesByBytes([]byte(testAliasesConfig))
	s.Nil(err)
	err = a.validate()
	s.Nil(err)
	invalid := `
packages:
  fasthttp: github.com/valyala/fasthttp
  fastrouter: github.com/fasthttp/router
bundles:
  webserver:
    - fasthttp
    - fastrouter
    - bitmap
`
	a, err = getAliasesByBytes([]byte(invalid))
	s.Nil(err)
	err = a.validate()
	s.NotNil(err)
}

func (s *modAliasesTestSuite) TestGlue() {
	a, err := getAliasesByBytes([]byte(testAliasesConfig))
	s.Nil(err)
	add := `
packages:
  fasthttp: github.com/valyala/fasthttp@v1.20.0
  bitmap: github.com/satmaelstorm/bitmap
  rtree: github.com/tidwall/rtree
bundles:
  search:
    - bitmap
    - rtree
`
	aa, err := getAliasesByBytes([]byte(add))
	s.Nil(err)
	a, err = a.Glue(aa)
	s.Nil(err)
	s.Contains(a.Packages, "bitmap")
	s.Contains(a.Packages, "rtree")
	s.Contains(a.Packages, "fasthttp")
	s.Contains(a.Packages, "fastrouter")
	s.Contains(a.Bundles, "webserver")
	s.Contains(a.Bundles, "search")
	s.Equal(a.Packages["fasthttp"], "github.com/valyala/fasthttp@v1.20.0")
	s.Equal(2, len(a.Bundles["search"]))
	s.Equal("bitmap", string(a.Bundles["search"][0]))
	s.Equal("rtree", string(a.Bundles["search"][1]))
}

func (s *modAliasesTestSuite) TestGetDefAliases() {
	a := GetAliases()
	s.True(len(a.Packages) > 0)
	s.True(len(a.Bundles) > 0)
	str := GetDefaultAliasesHelp()
	s.True(len(str) > 0)
}
