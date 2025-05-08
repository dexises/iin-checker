package databases

import (
	"context"
	"github.com/dexises/iin-checker/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func Connect(cfg config.DatabaseConfiguration) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(cfg.DSURL)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	return client.Database(cfg.DSDB), nil
}
