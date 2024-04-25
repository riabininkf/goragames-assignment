package repository

import (
	"errors"
	"github.com/riabininkf/goragames-assignment/internal/config"
	"github.com/riabininkf/goragames-assignment/internal/container"
	"github.com/riabininkf/goragames-assignment/internal/db"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const DefCodesName = "repository.codes"

func init() {
	container.Add(di.Def{
		Name: DefCodesName,
		Build: func(ctn di.Container) (interface{}, error) {
			var client *mongo.Client
			if err := container.Fill(db.DefMongoName, &client); err != nil {
				return nil, err
			}

			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return nil, err
			}

			var database string
			if database = cfg.GetString("mongodb.database"); database == "" {
				return nil, errors.New("mongodb.database config key is required")
			}

			return NewCodes(client.Database(database).Collection("codes")), nil
		},
	})
}
