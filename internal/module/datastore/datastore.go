package datastore

import (
	"context"
	"crud/internal/domain"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package datastore_test

type (
	mongodb interface {
		FindOne(ctx context.Context, filter interface{},
			opts ...*options.FindOneOptions) *mongo.SingleResult
		Find(ctx context.Context, filter interface{},
			opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
		DeleteOne(ctx context.Context, filter interface{},
			opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		UpdateOne(ctx context.Context, filter interface{}, update interface{},
			opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		InsertOne(ctx context.Context, document interface{},
			opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	}
	Datastore struct {
		l   zap.Logger
		co  mongodb
		ctx context.Context
	}
)

func NewDatastore(l zap.Logger, client mongodb, ctx context.Context) *Datastore {
	return &Datastore{
		l:   l,
		co:  client,
		ctx: ctx,
	}
}

func (ds *Datastore) AddUser(user domain.User) error {
	_, err := ds.co.InsertOne(ds.ctx, user)
	return err
}

func (ds *Datastore) GetUser(id string) (*domain.User, error) {
	var result domain.User
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err := ds.co.FindOne(ds.ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (ds *Datastore) GetList() ([]*domain.User, error) {
	users := []*domain.User{}
	cursor, err := ds.co.Find(ds.ctx, bson.M{})
	if err != nil || cursor.Err() != nil {
		if cursor.Err() != nil {
			return nil, cursor.Err()
		}
		return nil, err
	}
	defer cursor.Close(ds.ctx)
	for cursor.Next(ds.ctx) {
		var user *domain.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ds *Datastore) DeleteUser(id string) error {
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	res, err := ds.co.DeleteOne(ds.ctx, filter)
	if res.DeletedCount == 0 {
		return fmt.Errorf("no document to delete")
	}
	return err
}

func (ds *Datastore) UpdateUser(id string, newdata domain.User) error {

	newdata.ID = &id
	pByte, err := bson.Marshal(newdata)
	if err != nil {
		return err
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	res, err := ds.co.UpdateOne(ds.ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no document to update found")
	}
	return nil
}
