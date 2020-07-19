package application

import (
	"testing"

	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/domain/mocks"
	domain_mocks "github.com/concepts-system/go-paperless/domain/mocks"

	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	Config           *config.Configuration
	TokenKeyResolver TokenKeyResolver
	UsersMock        *domain_mocks.Users
	PasswordHelper   *passwordHelper
}

func (s *serviceTestSuite) SetupTest() {
	s.Config = &config.Configuration{}
	s.TokenKeyResolver = ConfigTokenKeyResolver(s.Config)
	s.UsersMock = &mocks.Users{}
	s.PasswordHelper = &passwordHelper{}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &serviceTestSuite{})
}
