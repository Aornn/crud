package datastore_test

import (
	"bufio"
	"bytes"
	"context"
	"crud/internal/domain"
	"crud/internal/module/datastore"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func TestAddUser(t *testing.T) {
	user := domain.User{}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockmongodb(ctrl)
	_, l, _ := logger()
	ds := datastore.NewDatastore(*l, db, ctx)
	t.Run("Valid", func(t *testing.T) {
		db.EXPECT().InsertOne(ctx, user).Return(nil, nil)
		err := ds.AddUser(user)
		assert.Nil(t, err)
	})
	t.Run("InsertOne returns error", func(t *testing.T) {
		db.EXPECT().InsertOne(ctx, user).Return(nil, errors.New("some error"))
		err := ds.AddUser(user)
		assert.EqualError(t, err, "some error")
	})
}

func TestGetUser(t *testing.T) {
	id := "1"
	user := domain.User{}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockmongodb(ctrl)
	_, l, _ := logger()
	ds := datastore.NewDatastore(*l, db, ctx)
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	t.Run("Valid", func(t *testing.T) {
		db.EXPECT().FindOne(ctx, filter).Return(mongo.NewSingleResultFromDocument(user, nil, nil))
		out, err := ds.GetUser(id)
		assert.Equal(t, out, &user)
		assert.Nil(t, err)
	})
	t.Run("InsertOne returns error", func(t *testing.T) {
		db.EXPECT().FindOne(ctx, filter).Return(mongo.NewSingleResultFromDocument(user, errors.New("error"), nil))
		u, err := ds.GetUser(id)
		assert.Nil(t, u)
		assert.EqualError(t, err, "error")
	})
}

func TestGetList(t *testing.T) {
	type dataS struct {
		Users []interface{}
	}

	data := dataS{}
	data.Users = append(data.Users, domain.User{})
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockmongodb(ctrl)
	_, l, _ := logger()
	ds := datastore.NewDatastore(*l, db, ctx)
	t.Run("Valid", func(t *testing.T) {
		db.EXPECT().Find(ctx, bson.M{}).Return(mongo.NewCursorFromDocuments(data.Users, nil, nil))
		out, err := ds.GetList()
		assert.Equal(t, out[0], &domain.User{})
		assert.Nil(t, err)
	})
	t.Run("Find returns error", func(t *testing.T) {
		db.EXPECT().Find(ctx, bson.M{}).Return(mongo.NewCursorFromDocuments(data.Users, errors.New("errors"), nil))
		out, err := ds.GetList()
		assert.Nil(t, out)
		assert.EqualError(t, err, "errors")
	})
}

func TestDeleteUser(t *testing.T) {
	id := "1"
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockmongodb(ctrl)
	_, l, _ := logger()
	ds := datastore.NewDatastore(*l, db, ctx)
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	t.Run("Valid", func(t *testing.T) {
		db.EXPECT().DeleteOne(ctx, filter).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
		err := ds.DeleteUser(id)
		assert.Nil(t, err)
	})
	t.Run("Valid no count", func(t *testing.T) {
		db.EXPECT().DeleteOne(ctx, filter).Return(&mongo.DeleteResult{DeletedCount: 0}, nil)
		err := ds.DeleteUser(id)
		assert.EqualError(t, err, "no document to delete")
	})
	t.Run("DeleteOne returns error", func(t *testing.T) {
		db.EXPECT().DeleteOne(ctx, filter).Return(&mongo.DeleteResult{DeletedCount: 1}, errors.New("error"))
		err := ds.DeleteUser(id)
		assert.EqualError(t, err, "error")
	})
}

func TestUpdateUser(t *testing.T) {
	id := "1"
	user := domain.User{}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockmongodb(ctrl)
	_, l, _ := logger()
	ds := datastore.NewDatastore(*l, db, ctx)
	filter := bson.D{primitive.E{Key: "id", Value: id}}

	t.Run("Valid", func(t *testing.T) {
		db.EXPECT().UpdateOne(ctx, filter, gomock.Any()).Return(&mongo.UpdateResult{MatchedCount: 1}, nil)
		err := ds.UpdateUser(id, user)
		assert.Nil(t, err)
	})
	t.Run("UpdateOne returns error", func(t *testing.T) {
		db.EXPECT().UpdateOne(ctx, filter, gomock.Any()).Return(&mongo.UpdateResult{MatchedCount: 1}, errors.New("some error"))
		err := ds.UpdateUser(id, user)
		assert.EqualError(t, err, "some error")
	})
	t.Run("UpdateOne returns no update", func(t *testing.T) {
		db.EXPECT().UpdateOne(ctx, filter, gomock.Any()).Return(&mongo.UpdateResult{MatchedCount: 0}, nil)
		err := ds.UpdateUser(id, user)
		assert.EqualError(t, err, "no document to update found")
	})
}
