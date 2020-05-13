package main

import (
	"context"
	"fmt"
	_movieHttpDelivery "github.com/cristhoperdev/go-mongo-clean/movie/delivery/http"
	_movieHttpDeliveryMiddleware "github.com/cristhoperdev/go-mongo-clean/movie/delivery/http/middleware"
	_movieRepository "github.com/cristhoperdev/go-mongo-clean/movie/repository/mongodb"
	_movieUseCase "github.com/cristhoperdev/go-mongo-clean/movie/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func init() {
	viper.SetConfigFile(`./config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.mongodb.host`)
	dbPort := viper.GetString(`database.mongodb.port`)
	dbase := viper.GetString(`database.mongodb.db`)

	ctx := context.TODO()

	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbase)
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _movieHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	movieRepository := _movieRepository.NewMongoMovieRepository(db, "movie")

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _movieUseCase.NewMovieUsecase(movieRepository, timeoutContext)
	_movieHttpDelivery.NewMovieHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address")))
}