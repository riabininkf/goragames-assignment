package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/riabininkf/goragames-assignment/internal/config"
	"github.com/riabininkf/goragames-assignment/internal/container"
	"github.com/riabininkf/goragames-assignment/internal/db"
	"github.com/riabininkf/goragames-assignment/internal/repository"
	"github.com/riabininkf/goragames-assignment/internal/repository/entity"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// todo find migration library that supports mongo
func init() {
	cmd := &cobra.Command{Use: "migrate"}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Apply all available migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			var client *mongo.Client
			if err := container.Fill(db.DefMongoName, &client); err != nil {
				return err
			}

			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return err
			}

			var tx repository.Tx
			if err := container.Fill(repository.DefTxName, &tx); err != nil {
				return err
			}

			var dbName string
			if dbName = cfg.GetString("mongodb.database"); dbName == "" {
				return errors.New("mongodb.database config key is required")
			}

			ctx, cancelFunc := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancelFunc()

			database := client.Database(dbName)

			if err := tx.Do(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
				if _, err := database.Collection("codes").Indexes().CreateOne(sessCtx, mongo.IndexModel{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetUnique(true),
				}); err != nil {
					return nil, fmt.Errorf("can't create indexes: %w", err)
				}

				if _, err := database.Collection("codes").InsertMany(sessCtx, []interface{}{
					entity.Code{Name: "twitter-reg1", Usages: 10},
					entity.Code{Name: "telegram-test", Usages: 10},
					entity.Code{Name: "instagram-hello", Usages: 10},
				}); err != nil {
					return nil, fmt.Errorf("can't insert codes: %w", err)
				}

				if _, err := database.Collection("candidates").Indexes().CreateOne(sessCtx, mongo.IndexModel{
					Keys:    bson.D{{Key: "email", Value: 1}, {Key: "code", Value: 1}},
					Options: options.Index().SetUnique(true),
				}); err != nil {
					return nil, fmt.Errorf("can't create indexes: %w", err)
				}

				return nil, nil
			}); err != nil {
				return err
			}

			return nil
		},
	})

	RootCmd.AddCommand(cmd)
}
