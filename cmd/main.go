package main

import (
	"context"
	"crud/internal/handler/add"
	"crud/internal/handler/addinit"
	"crud/internal/handler/delete"
	"crud/internal/handler/get"
	"crud/internal/handler/getlist"
	"crud/internal/handler/login"
	"crud/internal/handler/put"
	"crud/internal/module/datastore"
	"crud/internal/module/filestorer"
	usecasadd "crud/internal/usecase/add"
	usecasdel "crud/internal/usecase/del"
	usecasget "crud/internal/usecase/get"
	usecasgetlist "crud/internal/usecase/getlist"
	"fmt"
	"os"

	usecaselogin "crud/internal/usecase/login"
	usecaseupdate "crud/internal/usecase/put"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type conf struct {
	dsn           string
	db            string
	table         string
	inputfilepath string
	username      string
	password      string
}

func initConf() (*conf, error) {
	if os.Getenv("DSN") == "" {
		return nil, fmt.Errorf("no dsn")
	}
	if os.Getenv("DB") == "" {
		return nil, fmt.Errorf("no database")
	}
	if os.Getenv("USERNAME") == "" {
		return nil, fmt.Errorf("no dsn")
	}
	if os.Getenv("PASSWORD") == "" {
		return nil, fmt.Errorf("no database")
	}
	if os.Getenv("TABLE") == "" {
		return nil, fmt.Errorf("no table")
	}
	if os.Getenv("FILE") == "" {
		return nil, fmt.Errorf("no file")
	}
	return &conf{
		dsn:           os.Getenv("DSN"),
		db:            os.Getenv("DB"),
		table:         os.Getenv("TABLE"),
		inputfilepath: os.Getenv("FILE"),
		username:      os.Getenv("USERNAME"),
		password:      os.Getenv("PASSWORD"),
	}, nil
}

func main() {
	confData, err := initConf()
	if err != nil {
		panic(err)
	}
	l, err := zap.NewProduction(zap.WithCaller(false), zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		panic("unable to get logger: " + err.Error())
	}
	ctx := context.Background()
	cred := options.Credential{
		Username: confData.username,
		Password: confData.password,
	}
	clientOptions := options.Client().ApplyURI(confData.dsn).SetAuth(cred)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	ds := datastore.NewDatastore(*l, client.Database(confData.db).Collection(confData.table), ctx)
	fs := filestorer.NewFileStorer()
	procAdd := usecasadd.NewUsecase(ds, fs)
	procLogin := usecaselogin.NewUsecase(ds)
	handlerInit := addinit.New(*l, procAdd)
	err = handlerInit.Handle(confData.inputfilepath)
	if err != nil && err.Error() != "already present" {
		panic(err)
	}
	handlerLogin := login.New(*l, procLogin)

	procGet := usecasget.NewUsecase(ds)
	procGetList := usecasgetlist.NewUsecase(ds)

	procDel := usecasdel.NewUsecase(ds, fs)
	procUpdate := usecaseupdate.NewUsecase(ds, fs)
	handlerAdd := add.New(*l, procAdd)
	handlerGet := get.New(*l, procGet)
	handlerGetList := getlist.New(*l, procGetList)
	handlerPut := put.New(*l, procUpdate)
	handlerDel := delete.New(*l, procDel)
	router := gin.Default()

	router.POST("/login", handlerLogin.Handle)
	router.POST("/add/user", handlerAdd.Handle)
	router.DELETE("/delete/user/:id", handlerDel.Handle)
	router.GET("/user/list", handlerGetList.Handle)
	router.GET("/user/:id", handlerGet.Handle)
	router.PUT("/user/:id", handlerPut.Handle)

	router.Run()
}
