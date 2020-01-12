package application

import (
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
	"golang.org/x/crypto/bcrypt"
)

type passwordHelper struct{}

func (passwordHelper) checkUserPassword(user *domain.User, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
}

func (passwordHelper) setUserPassword(user *domain.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return errors.Wrapf(err, "Failed to hash password")
	}

	user.Password = domain.Password(string(hash))
	return nil
}
