package put

import (
	"crud/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package put_test

type (
	iUpdateUserFromDatabase interface {
		UpdateUser(id string, newdata domain.User) error
		GetUser(id string) (*domain.User, error)
	}
	iUpdateUserFile interface {
		UpdateUser(filename, data string) error
	}
	// Handler : struct which provide the variables across the program.
	Usecase struct {
		datastore  iUpdateUserFromDatabase
		filestorer iUpdateUserFile
	}
)

func NewUsecase(ds iUpdateUserFromDatabase, fs iUpdateUserFile) *Usecase {
	return &Usecase{
		datastore:  ds,
		filestorer: fs,
	}
}

func updateData(new *domain.User, old *domain.User) {
	if new.Password == nil {
		new.Password = old.Password
	}
	if new.IsActive == nil {
		new.IsActive = old.IsActive
	}
	if new.Balance == nil {
		new.Balance = old.Balance
	}
	if new.Age == nil {
		new.Age = old.Age
	}
	if new.Gender == nil {
		new.Gender = old.Gender
	}
	if new.Company == nil {
		new.Company = old.Company
	}
	if new.Email == nil {
		new.Email = old.Email
	}
	if new.Phone == nil {
		new.Phone = old.Phone
	}
	if new.Address == nil {
		new.Address = old.Address
	}
	if new.About == nil {
		new.About = old.About
	}
	if new.Registred == nil {
		new.Registred = old.Registred
	}
	if new.Latitude == nil {
		new.Latitude = old.Latitude
	}
	if new.Longitude == nil {
		new.Longitude = old.Longitude
	}
	if new.Tags == nil {
		new.Tags = old.Tags
	}
	if new.Friends == nil {
		new.Friends = old.Friends
	}
	if new.Data == nil {
		new.Friends = old.Friends
	}
}

func (u *Usecase) Process(l *zap.Logger, id string, newdata domain.User) error {
	var err error
	if newdata.Data == nil {
		return fmt.Errorf("can't continue without data")
	}
	if newdata.ID != nil && *newdata.ID != id {
		return fmt.Errorf("not same id")
	}
	if newdata.Password != nil {
		newdata.Password, err = domain.HashPassword(newdata.Password)
		if err != nil {
			return err
		}
	}
	uget, err := u.datastore.GetUser(id)
	if err != nil {
		return err
	}
	updateData(&newdata, uget)
	l.Info("updating", zap.String("id", id))
	err = u.datastore.UpdateUser(id, newdata)
	if err != nil {
		return err
	}
	return u.filestorer.UpdateUser(id, *newdata.Data)
}
