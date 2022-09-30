package delete_test

import (
	"bufio"
	"bytes"
	"crud/internal/handler/delete"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() (*bytes.Buffer, *zap.Logger, *bufio.Writer) {
	var buffer bytes.Buffer
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.TimeKey = ""
	encoder := zapcore.NewJSONEncoder(encoderConf)
	writer := bufio.NewWriter(&buffer)

	l := zap.New(zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(writer)), zapcore.DebugLevel))

	return &buffer, l, writer
}

func createHeader(c *gin.Context) {
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
}

func TestHandle(t *testing.T) {
	var ()
	t.Run("Valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := NewMockiDelUser(ctrl)
		_, l, _ := logger()
		h := delete.New(*l, p)
		createHeader(c)
		p.EXPECT().Process(l, c.Param("id")).Return(nil)
		h.Handle(c)
		assert.EqualValues(t, http.StatusOK, w.Code)
	})

	t.Run("process returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := NewMockiDelUser(ctrl)
		_, l, _ := logger()
		h := delete.New(*l, p)
		createHeader(c)
		p.EXPECT().Process(l, c.Param("id")).Return(errors.New("error"))
		h.Handle(c)
		assert.EqualValues(t, http.StatusBadRequest, w.Code)
	})
}
