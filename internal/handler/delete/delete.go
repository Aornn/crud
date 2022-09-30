package delete

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package delete_test

type (
	iDelUser interface {
		Process(l *zap.Logger, id string) error
	}
	// Handler : struct which provide the variables across the program.
	Handler struct {
		l zap.Logger
		p iDelUser
	}
)

// New : Return a pointer on new Handler.
func New(l zap.Logger, p iDelUser) *Handler {
	return &Handler{l: l, p: p}
}

func (h *Handler) Handle(c *gin.Context) {
	err := h.p.Process(&h.l, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": c.Param("id"), "status": "deleted"})

}
