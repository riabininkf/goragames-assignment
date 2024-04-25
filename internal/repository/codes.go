package repository

//go:generate mockery --name=Codes --inpackage --output=. --filename=codes_mock.go --structname=CodesMock

import (
	"context"
	"errors"
	"github.com/riabininkf/goragames-assignment/internal/repository/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewCodes(collection *mongo.Collection) Codes {
	return &codes{collection: collection}
}

type (
	Codes interface {
		GetByNameWithLock(ctx context.Context, name string) (*entity.Code, error)
		DecrementUsagesByName(ctx context.Context, name string) error
	}

	codes struct {
		client     *mongo.Client
		collection *mongo.Collection
	}
)

func (r *codes) GetByNameWithLock(ctx context.Context, name string) (*entity.Code, error) {
	var result entity.Code
	if err := r.collection.FindOneAndUpdate(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"lockedAtMs": time.Now().UnixMicro()}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *codes) DecrementUsagesByName(ctx context.Context, name string) error {
	if _, err := r.collection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$inc": bson.M{"usages": -1}}); err != nil {
		return err
	}

	return nil
}
