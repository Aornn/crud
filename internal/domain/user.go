package domain

import (
	"crud/pkg"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"golang.org/x/crypto/bcrypt"
)

type (
	Friend struct {
		ID   int
		Name string
	}
	User struct {
		ID        *string
		Password  *string
		IsActive  *bool
		Balance   *string
		Age       *string
		Gender    *string
		Company   *string
		Email     *string
		Phone     *string
		Address   *string
		About     *string
		Registred *string
		Latitude  *float64
		Longitude *float64
		Tags      *[]string
		Friends   *[]Friend
		Data      *string
	}
	UserLogin struct {
		ID       string
		Password string
	}
)

func HashPassword(input *string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	strHash := string(hashedPassword)
	return &strHash, nil
}

func ComparePassword(input *string, hash *string) error {
	return bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*input))
}

func LoginToDomain(user *pkg.UserLogin) *UserLogin {
	return &UserLogin{
		ID:       user.ID,
		Password: user.Password,
	}
}

func ToPkg(user *User) *pkg.User {

	out := pkg.User{
		ID:        user.ID,
		Password:  user.Password,
		IsActive:  user.IsActive,
		Balance:   user.Balance,
		Gender:    user.Gender,
		Company:   user.Company,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		About:     user.About,
		Registred: user.Registred,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
		Tags:      user.Tags,
		Friends:   nil,
		Data:      user.Data,
	}
	if user.Age == nil {
		out.Age = nil
	} else {
		age := json.Number(*user.Age)
		out.Age = &age
	}
	if user.Friends != nil {
		out.Friends = &[]pkg.Friend{}
		for _, e := range *user.Friends {
			*out.Friends = append(*out.Friends, pkg.Friend{
				ID:   e.ID,
				Name: e.Name,
			})
		}
	}
	return &out
}

func ToDomain(user pkg.User) User {
	out := User{
		ID:        user.ID,
		Password:  user.Password,
		IsActive:  user.IsActive,
		Balance:   user.Balance,
		Gender:    user.Gender,
		Company:   user.Company,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		About:     user.About,
		Registred: user.Registred,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
		Tags:      user.Tags,
		Friends:   nil,
		Data:      user.Data,
	}
	if user.Age != nil {
		out.Age = aws.String(user.Age.String())
	} else {
		out.Age = nil
	}
	if user.Friends != nil {
		out.Friends = &[]Friend{}
		for _, e := range *user.Friends {
			*out.Friends = append(*out.Friends, Friend{
				ID:   e.ID,
				Name: e.Name,
			})
		}
	}

	return out
}
