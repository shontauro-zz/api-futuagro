package store

import (
	"context"
	"fmt"
	"time"

	"futuagro.com/pkg/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDB return a mongodb connection
func NewDB(confPtr *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(confPtr.Database.URI)
	clientOptions.SetMaxPoolSize(confPtr.Database.PoolSize)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating a mongoDB client")
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Error connecting to MongoDB")
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrapf(err, "Error connecting (Ping) to MongoDB")
	}

	fmt.Printf("Connected to database MongoDB :%s \n", confPtr.Database.URI)
	return client, nil
}
