package addinit_test

import (
	"bufio"
	"bytes"
	"crud/internal/domain"
	"crud/internal/handler/addinit"
	"crud/pkg"
	"errors"
	"os"
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

func TestHandle(t *testing.T) {
	var (
		inputfilename = "test"
		user          = pkg.User{ID: aws.String("1")}
		domainU       = domain.ToDomain(user)
	)
	t.Run("file does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, l, _ := logger()
		p := NewMockiAddUser(ctrl)
		h := addinit.New(*l, p)
		err := h.Handle(inputfilename)
		assert.EqualError(t, err, "open test: no such file or directory")
	})
	t.Run("can't unmarshal", func(t *testing.T) {
		fs, err := os.Create(inputfilename)
		if err != nil {
			panic(err)
		}

		_, err = fs.WriteString("aze")
		if err != nil {
			panic(err)
		}
		fs.Close()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, l, _ := logger()
		p := NewMockiAddUser(ctrl)
		h := addinit.New(*l, p)
		err = h.Handle(inputfilename)
		assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
		err = os.Remove(inputfilename)
		if err != nil {
			panic(err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		fs, err := os.Create(inputfilename)
		if err != nil {
			panic(err)
		}

		_, err = fs.WriteString(`[{"id": "1"}]`)
		if err != nil {
			panic(err)
		}
		fs.Close()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, l, _ := logger()
		p := NewMockiAddUser(ctrl)
		h := addinit.New(*l, p)
		p.EXPECT().Process(l, &domainU).Times(1).Return(nil)
		err = h.Handle(inputfilename)
		assert.Nil(t, err)
		err = os.Remove(inputfilename)
		if err != nil {
			panic(err)
		}
	})

	t.Run("process returns error", func(t *testing.T) {
		fs, err := os.Create(inputfilename)
		if err != nil {
			panic(err)
		}

		_, err = fs.WriteString(`[{"id": "1"}]`)
		if err != nil {
			panic(err)
		}
		fs.Close()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, l, _ := logger()
		p := NewMockiAddUser(ctrl)
		h := addinit.New(*l, p)
		p.EXPECT().Process(l, &domainU).Times(1).Return(errors.New("error"))
		err = h.Handle(inputfilename)
		assert.EqualError(t, err, "error")
		err = os.Remove(inputfilename)
		if err != nil {
			panic(err)
		}
	})
}
