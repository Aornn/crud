package put_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/handler/put"
	"crud/pkg"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

func createjson(c *gin.Context, input *pkg.User) {
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	jsonbytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: "1"},
	}

}

func TestHandle(t *testing.T) {
	var (
		input       = pkg.User{ID: aws.String("1")}
		domainInput = domain.ToDomain(input)
	)
	t.Run("Valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := NewMockiPutUser(ctrl)
		_, l, _ := logger()
		h := put.New(*l, p)
		createjson(c, &input)
		p.EXPECT().Process(l, c.Param("id"), domainInput).Return(nil)
		h.Handle(c)
		assert.EqualValues(t, http.StatusOK, w.Code)
	})
	t.Run("Process returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := NewMockiPutUser(ctrl)
		_, l, _ := logger()
		h := put.New(*l, p)
		createjson(c, &input)
		p.EXPECT().Process(l, c.Param("id"), domainInput).Return(errors.New("some error"))
		h.Handle(c)
		assert.EqualValues(t, http.StatusBadRequest, w.Code)
	})
	t.Run("ShouldBindJSONreturns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := NewMockiPutUser(ctrl)
		_, l, _ := logger()
		h := put.New(*l, p)
		createjson(c, nil)
		c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte("sfdsf")))
		h.Handle(c)
		assert.EqualValues(t, http.StatusBadRequest, w.Code)
	})
}
