package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/riabininkf/goragames-assignment/internal/config"
	"github.com/riabininkf/goragames-assignment/internal/container"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	DefMongoName = "db.mongo"

	defaultTTL = time.Second * 5
)

func init() {
	container.Add(di.Def{
		Name: DefMongoName,
		Build: func(ctn di.Container) (interface{}, error) {
			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return nil, err
			}

			var connURI string
			if connURI = cfg.GetString("mongodb.conn"); connURI == "" {
				return nil, errors.New("mongodb.conn config key is required")
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTTL)
			defer cancelFunc()

			var (
				err    error
				client *mongo.Client
			)
			if client, err = mongo.Connect(ctx, options.Client().ApplyURI(connURI)); err != nil {
				return nil, fmt.Errorf("can't connect to mongo, %w", err)
			}

			if err = client.Ping(ctx, nil); err != nil {
				return nil, fmt.Errorf("can't ping mongo, %w", err)
			}

			return client, nil
		},
		Close: func(obj interface{}) error {
			ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTTL)
			defer cancelFunc()

			_ = obj.(*mongo.Client).Disconnect(ctx)
			return nil
		},
	})
}
