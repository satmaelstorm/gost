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


