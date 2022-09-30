package getlist

import (
	"crud/internal/domain"
	"crud/pkg"
	"sync"

	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package getlist_test

type (
	iGetUserFromDatabase interface {
		GetList() ([]*domain.User, error)
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

func convert(users *[]*pkg.User, user *domain.User, wg *sync.WaitGroup) {
	defer wg.Done()
	toAdd := domain.ToPkg(user)
	*users = append(*users, toAdd)
}

func (u *Usecase) Process(l *zap.Logger) ([]*pkg.User, error) {
	var wg sync.WaitGroup
	l.Info("searching for all users")
	output := []*pkg.User{}
	users, err := u.datastore.GetList()
	if err != nil {
		return nil, err
	}
	for _, e := range users {
		wg.Add(1)
		go convert(&output, e, &wg)
	}
	wg.Wait()
	return output, nil
}
