package addinit

import (
	"context"
	"crud/internal/domain"
	"crud/pkg"
	"encoding/json"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package addinit_test

type (
	iAddUser interface {
		Process(l *zap.Logger, user *domain.User) error
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iAddUser
	}
)

func New(l zap.Logger, p iAddUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(inputfile string) error {
	errs, _ := errgroup.WithContext(context.Background())

	var users []pkg.User
	jsonFile, err := os.Open(inputfile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteResult, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteResult, &users)
	if err != nil {
		h.l.Info("can't unmarshal", zap.Error(err))
		return err
	}
	for _, u := range users {
		user := domain.ToDomain(u)
		errs.Go(func() error {
			return h.p.Process(&h.l, &user)
		})
	}
	return errs.Wait()
}
