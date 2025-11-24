package database

import (
	"context"
	"time"
	"log"

	"prestasi_backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func ConnectMongo() (*mongo.Database, error) {
	uri := config.Get("MONGO_URI")
	dbName := config.Get("MONGO_DB")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	log.Println("✅ MongoDB berhasil terkoneksi")
	return client.Database(dbName), nil
}
