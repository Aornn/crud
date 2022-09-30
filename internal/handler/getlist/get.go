package getlist

import (
	"crud/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package getlist_test

type (
	iGetList interface {
		Process(l *zap.Logger) ([]*pkg.User, error)
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iGetList
	}
)

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iGetList) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	users, err := h.p.Process(&h.l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
