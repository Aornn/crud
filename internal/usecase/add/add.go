package add

import (
	"crud/internal/domain"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package add_test

type (
	iAddUserInDatabase interface {
		AddUser(user domain.User) error
		GetUser(id string) (*domain.User, error)
	}
	iWriteUserToFile interface {
		WriteUser(filename, data string) error
		Checkfile(filename string) bool
	}
	// Handler : struct which provide the variables across the program.
	Usecase struct {
		datastore  iAddUserInDatabase
		filestorer iWriteUserToFile
	}
)

func NewUsecase(ds iAddUserInDatabase, fs iWriteUserToFile) *Usecase {
	return &Usecase{
		datastore:  ds,
		filestorer: fs,
	}
}

func (u *Usecase) Process(l *zap.Logger, user *domain.User) error {
	var err error
	if user.Data == nil {
		return fmt.Errorf("can't continue without data")
	}
	if user.ID == nil {
		return fmt.Errorf("can't continue without id")
	}
	if user.Password != nil {
		user.Password, err = domain.HashPassword(user.Password)
		if err != nil {
			return err
		}
	}
	uget, err := u.datastore.GetUser(*user.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	} else if uget != nil {
		l.Info("already present in database", zap.String("id", *user.ID))
		if !u.filestorer.Checkfile(*user.ID) {
			l.Info("but file not found")
			err = u.filestorer.WriteUser(*user.ID, *user.Data)
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("already present")
	}
	err = u.filestorer.WriteUser(*user.ID, *user.Data)
	if err != nil {
		return err
	}
	err = u.datastore.AddUser(*user)
	if err != nil {
		return err
	}
	return nil
}
