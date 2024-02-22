package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Config struct {
	MongoClient *mongo.Client
	MemoryDB    map[string]string
	Port        string
}

func NewConfig(dbUri string) *Config {
	dbClient := connectDB(dbUri)
	memoryMap := make(map[string]string)
	return &Config{
		MongoClient: dbClient,
		MemoryDB:    memoryMap,
		Port:        os.Getenv("PORT"),
	}
}

func connectDB(dbUri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatalf("MongoDB istemcisi oluşturulurken hata: %v", err)
	}

	// Creating a context to control connecting to DB. If db can not connect is 30 seconds
	// cancel is called and connection is closed.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connection database: %v", err)
	}

	if client != nil {
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Error trying ping database: %v", err)
		}
	} else {
		log.Fatal("Mongo client is nil, can not check connection")
	}

	log.Println("Pinged database!")
	return client
}
