package login_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/usecase/login"
	"errors"
	"testing"

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
func TestProcess(t *testing.T) {
	var (
		id      = "1"
		hashpwd = "$2a$10$7GifBZUH1FSCxPQpZsX8TOfJtBD6A6vZ6tHG88OazlHrtkBfCaqjK"
		ret     = domain.User{
			ID:       &id,
			Password: &hashpwd,
		}
		log = domain.UserLogin{
			ID:       id,
			Password: "123",
		}
	)
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiGetUserFromDatabase(ctrl)
		p := login.NewUsecase(ds)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(&ret, nil)
		u, err := p.Process(l, &log)
		assert.Nil(t, err)
		assert.Equal(t, u, &ret)
	})

	t.Run("GetUser returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiGetUserFromDatabase(ctrl)
		p := login.NewUsecase(ds)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(nil, errors.New("error"))
		u, err := p.Process(l, &log)
		assert.Nil(t, u)
		assert.EqualError(t, err, "error")
	})
}
