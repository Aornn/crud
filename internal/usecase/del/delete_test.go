package del_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/usecase/del"
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
		user = domain.User{
			ID: &id,
		}
	)
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiDeleteUserFromDatabase(ctrl)
		fs := NewMockiDeleteUserFile(ctrl)
		p := del.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().DeleteUser(*user.ID).Return(nil)
		fs.EXPECT().DeleteUser(*user.ID).Return(nil)
		err := p.Process(l, *user.ID)
		assert.Nil(t, err)
	})

	t.Run("ds.DeleteUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiDeleteUserFromDatabase(ctrl)
		fs := NewMockiDeleteUserFile(ctrl)
		p := del.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().DeleteUser(*user.ID).Return(errors.New("some error"))
		err := p.Process(l, *user.ID)
		assert.EqualError(t, err, "some error")
	})

	t.Run("fs.DeleteUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiDeleteUserFromDatabase(ctrl)
		fs := NewMockiDeleteUserFile(ctrl)
		p := del.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().DeleteUser(*user.ID).Return(nil)
		fs.EXPECT().DeleteUser(*user.ID).Return(errors.New("error"))
		err := p.Process(l, *user.ID)
		assert.EqualError(t, err, "error")
	})

}
