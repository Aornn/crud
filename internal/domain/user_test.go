package domain_test

import (
	"crud/internal/domain"
	"crud/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	var (
		id   = "1"
		data = "data"
		user = domain.User{
			ID:   &id,
			Data: &data,
		}
		userLogin = domain.UserLogin{
			ID:       "123",
			Password: "123",
		}
		pkgUser = pkg.User{
			ID:   &id,
			Data: &data,
		}
		pkgLogin = pkg.UserLogin{
			ID:       "123",
			Password: "123",
		}
	)
	t.Run("ToDomain", func(t *testing.T) {
		assert.Equal(t, domain.ToDomain(pkgUser), user)
	})
	t.Run("ToPkg", func(t *testing.T) {
		assert.Equal(t, domain.ToPkg(&user), &pkgUser)
	})

	t.Run("UserLogin ToDomain", func(t *testing.T) {
		assert.Equal(t, domain.LoginToDomain(&pkgLogin), &userLogin)
	})
}
