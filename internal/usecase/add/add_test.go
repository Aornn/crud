package add_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/usecase/add"
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
		data = "data"
		user = domain.User{
			ID:   &id,
			Data: &data,
		}
		newUser = domain.User{}
	)
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(*user.ID).Return(nil, nil)
		fs.EXPECT().WriteUser(*user.ID, *user.Data).Return(nil)
		ds.EXPECT().AddUser(user).Return(nil)
		err := p.Process(l, &user)
		assert.Nil(t, err)
	})
	t.Run("In db and file not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(*user.ID).Return(&newUser, nil)
		fs.EXPECT().Checkfile(*user.ID).Return(false)
		fs.EXPECT().WriteUser(*user.ID, *user.Data).Return(nil)
		err := p.Process(l, &user)
		assert.EqualError(t, err, "already present")
	})
	t.Run("AddUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(*user.ID).Return(nil, nil)
		fs.EXPECT().WriteUser(*user.ID, *user.Data).Return(nil)
		ds.EXPECT().AddUser(user).Return(errors.New("error"))
		err := p.Process(l, &user)
		assert.EqualError(t, err, "error")
	})
	t.Run("WriteUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(*user.ID).Return(nil, nil)
		fs.EXPECT().WriteUser(*user.ID, *user.Data).Return(errors.New("error"))
		err := p.Process(l, &user)
		assert.EqualError(t, err, "error")
	})
	t.Run("GetUser returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		ds.EXPECT().GetUser(*user.ID).Return(nil, errors.New("error"))
		err := p.Process(l, &user)
		assert.EqualError(t, err, "error")
	})
	t.Run("user id is nil", func(t *testing.T) {
		userdatanil := domain.User{
			Data: &data,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		err := p.Process(l, &userdatanil)
		assert.EqualError(t, err, "can't continue without id")
	})

	t.Run("user data is nil", func(t *testing.T) {
		userdatanil := domain.User{
			ID: &id,
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ds := NewMockiAddUserInDatabase(ctrl)
		fs := NewMockiWriteUserToFile(ctrl)
		p := add.NewUsecase(ds, fs)
		_, l, _ := logger()
		err := p.Process(l, &userdatanil)
		assert.EqualError(t, err, "can't continue without data")
	})

}
