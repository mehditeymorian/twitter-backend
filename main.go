package main

import (
	"context"
	"github.com/arman-aminian/twitter-backend/db"
	_ "github.com/arman-aminian/twitter-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/arman-aminian/twitter-backend/handler"
	"github.com/arman-aminian/twitter-backend/router"
	"github.com/arman-aminian/twitter-backend/store"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	r := router.New()

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	mongoClient, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	usersDb := mongoClient.Database("twitter_db").Collection("users")
	createUniqueIndices(usersDb, "username")
	createUniqueIndices(usersDb, "email")

	g := r.Group("")

	us := store.NewUserStore(usersDb)

	h := handler.NewHandler(us)
	h.Register(g)

	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}

func createUniqueIndices(db *mongo.Collection, field string) {
	_, err := db.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
