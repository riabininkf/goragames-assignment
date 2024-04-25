package repository

import (
	"github.com/riabininkf/goragames-assignment/internal/container"
	"github.com/riabininkf/goragames-assignment/internal/db"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const DefTxName = "repository.tx"

func init() {
	container.Add(di.Def{
		Name: DefTxName,
		Build: func(ctn di.Container) (interface{}, error) {
			var client *mongo.Client
			if err := container.Fill(db.DefMongoName, &client); err != nil {
				return nil, err
			}

			return NewTx(client), nil
		},
	})
}
