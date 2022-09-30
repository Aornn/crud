package get

import (
	"crud/internal/domain"
	"crud/pkg"

	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package get_test

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

func (u *Usecase) Process(l *zap.Logger, id string) (*pkg.User, error) {
	l.Info("searching", zap.String("id", id))
	user, err := u.datastore.GetUser(id)
	if err != nil {
		return nil, err
	}
	return domain.ToPkg(user), nil
}
