package application

import (
	"errors"
	"testing"

	"github.com/concepts-system/go-paperless/domain"
)

const (
	testUsername     = "user"
	testUserPassword = "password"
	wrongUsername    = "nobody"
)

func testUser(suite *serviceTestSuite) *domain.User {
	user := domain.NewUser(domain.User{
		Username: testUsername,
		Forename: "Test",
		Surname:  "User",
		IsActive: true,
		IsAdmin:  false,
	})

	suite.PasswordHelper.setUserPassword(user, testUserPassword)
	return user
}

func (s *serviceTestSuite) TestAuthenticateUserByCredentials_WithCorrectCredentials() {
	user := testUser(s)
	s.UsersMock.On("GetByUsername", domain.Name(testUsername)).Return(user, nil)
	s.UsersMock.On("GetByUsername", domain.Name(wrongUsername)).Return(nil, errors.New("No such user"))

	cases := []struct {
		name          string
		username      string
		password      string
		expectSuccess bool
	}{
		{
			"CorrectCredentials",
			testUsername,
			testUserPassword,
			true,
		},
		{
			"InvalidPassword",
			testUsername,
			"wrongPassword",
			false,
		},
		{
			"InvalidUsername",
			wrongUsername,
			testUserPassword,
			false,
		},
	}

	service := NewAuthService(s.Config, s.UsersMock, s.TokenKeyResolver)

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			token, err := service.AuthenticateUserByCredentials(c.username, c.password)

			if c.expectSuccess {
				s.Assert().NotNil(token, "expected valid token")
				s.Assert().Nil(err, "did not expect an error")
				s.Assert().Equal(testUsername, token.Username, "did not expect an error")
			} else {
				s.Assert().Nil(token, "did not expect a valid token")
				s.Assert().NotNil(err, "did expect an error")
			}
		})
	}
}
