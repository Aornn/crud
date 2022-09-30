package getlist_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/usecase/getlist"
	"crud/pkg"
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
		id   = "1"
		user = []*pkg.User{
			{
				ID: &id,
			},
		}
		ret = []*domain.User{
			{
				ID: &id,
			},
		}
	)
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiGetUserFromDatabase(ctrl)
		p := getlist.NewUsecase(ds)
		_, l, _ := logger()
		ds.EXPECT().GetList().Return(ret, nil)
		u, err := p.Process(l)
		assert.Equal(t, u, user)
		assert.Nil(t, err)
	})

	t.Run("GetList returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiGetUserFromDatabase(ctrl)
		p := getlist.NewUsecase(ds)
		_, l, _ := logger()
		ds.EXPECT().GetList().Return(nil, errors.New("error"))
		u, err := p.Process(l)
		assert.EqualError(t, err, "error")
		assert.Nil(t, u)
	})
}
