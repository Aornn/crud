package get

import (
	"crud/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package get_test

type (
	iGetUser interface {
		Process(l *zap.Logger, id string) (*pkg.User, error)
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iGetUser
	}
)

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iGetUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	user, err := h.p.Process(&h.l, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
