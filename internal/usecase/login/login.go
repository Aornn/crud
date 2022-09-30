package login

import (
	"crud/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package login_test

type (
	iGetUserFromDatabase interface {
		GetUser(id string) (*domain.User, error)
	}

	// Handler : struct which provide the variables across the program.
	Usecase struct {
		datastore iGetUserFromDatabase
	}
)

func NewUsecase(ds iGetUserFromDatabase) *Usecase {
	return &Usecase{
		datastore: ds,
	}
}

func (u *Usecase) Process(l *zap.Logger, userLogin *domain.UserLogin) (*domain.User, error) {
	user, err := u.datastore.GetUser(userLogin.ID)
	if err != nil {
		return nil, err
	}
	if domain.ComparePassword(&userLogin.Password, user.Password) != nil {
		return nil, fmt.Errorf("wrong password")
	}

	return user, nil
}
