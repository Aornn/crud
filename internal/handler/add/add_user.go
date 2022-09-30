package add

import (
	"crud/internal/domain"
	"crud/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package add_test

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

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iAddUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	var user pkg.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	domainUser := domain.ToDomain(user)
	if err := h.p.Process(&h.l, &domainUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "you are stored"})
}
