package put

import (
	"crud/internal/domain"
	"crud/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package put_test

type (
	iPutUser interface {
		Process(l *zap.Logger, id string, newdata domain.User) error
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iPutUser
	}
)

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iPutUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	var user pkg.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	domainUser := domain.ToDomain(user)
	err := h.p.Process(&h.l, c.Param("id"), domainUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": c.Param("id"), "status": "updated"})

}
