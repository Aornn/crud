package del

import (
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package del_test

type (
	iDeleteUserFromDatabase interface {
		DeleteUser(id string) error
	}
	iDeleteUserFile interface {
		DeleteUser(id string) error
	}
	// Handler : struct which provide the variables across the program.
	Usecase struct {
		datastore  iDeleteUserFromDatabase
		filestorer iDeleteUserFile
	}
)

func NewUsecase(ds iDeleteUserFromDatabase, fs iDeleteUserFile) *Usecase {
	return &Usecase{
		datastore:  ds,
		filestorer: fs,
	}
}

func (u *Usecase) Process(l *zap.Logger, id string) error {
	l.Info("deleting", zap.String("id", id))
	err := u.datastore.DeleteUser(id)
	if err != nil {
		return err
	}
	return u.filestorer.DeleteUser(id)
}
