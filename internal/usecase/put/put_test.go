package put_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/usecase/put"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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
		id         = "1"
		newAddress = "Paris"
		data       = "data"
		newdata    = domain.User{
			Address: &newAddress,
			Data:    &data,
		}
		uget = domain.User{}
	)
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(&uget, nil)
		ds.EXPECT().UpdateUser(id, newdata).Return(nil)
		fs.EXPECT().UpdateUser(id, *newdata.Data).Return(nil)
		err := p.Process(l, id, newdata)
		assert.Nil(t, err)
	})

	t.Run("GetUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(nil, errors.New("error"))
		err := p.Process(l, id, newdata)
		assert.EqualError(t, err, "error")
	})

	t.Run("fs.UpdateUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(&uget, nil)
		ds.EXPECT().UpdateUser(id, newdata).Return(nil)
		fs.EXPECT().UpdateUser(id, *newdata.Data).Return(errors.New("error"))
		err := p.Process(l, id, newdata)
		assert.EqualError(t, err, "error")
	})

	t.Run("ds.UpdateUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(id).Return(&uget, nil)
		ds.EXPECT().UpdateUser(id, newdata).Return(errors.New("error"))
		err := p.Process(l, id, newdata)
		assert.EqualError(t, err, "error")
	})

	t.Run("no data", func(t *testing.T) {
		newdatanil := domain.User{}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		err := p.Process(l, id, newdatanil)
		assert.EqualError(t, err, "can't continue without data")
	})

	t.Run("id different", func(t *testing.T) {
		newdataid := domain.User{
			ID:   aws.String("9"),
			Data: &data,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiUpdateUserFromDatabase(ctrl)
		fs := NewMockiUpdateUserFile(ctrl)
		p := put.NewUsecase(ds, fs)
		_, l, _ := logger()
		err := p.Process(l, id, newdataid)
		assert.EqualError(t, err, "not same id")
	})

}
