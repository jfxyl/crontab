package mongo

import (
	"context"
	"crontab/internal/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() (err error) {
	global.MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(global.Config.MongoConfig.URI))
	if err != nil {
		return
	}
	global.LogCollection = global.MongoClient.Database("cron").Collection("log")
	return
}
