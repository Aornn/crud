package login

import (
	"crud/internal/domain"
	"crud/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package login_test

type (
	iLogUser interface {
		Process(l *zap.Logger, userLogin *domain.UserLogin) (*domain.User, error)
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iLogUser
	}
)

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iLogUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	var user pkg.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	domainUser := domain.LoginToDomain(&user)

	data, err := h.p.Process(&h.l, domainUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "logged", "user": domain.ToPkg(data)})

}
