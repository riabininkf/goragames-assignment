package repository

//go:generate mockery --name=Candidates --inpackage --output=. --filename=candidates_mock.go --structname=CandidatesMock

import (
	"context"
	"errors"
	"github.com/riabininkf/goragames-assignment/internal/repository/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCandidates(collection *mongo.Collection) Candidates {
	return &candidates{collection: collection}
}

type (
	Candidates interface {
		Add(ctx context.Context, candidate *entity.Candidate) error
		GetByEmailAndCode(ctx context.Context, email, code string) (*entity.Candidate, error)
	}

	candidates struct {
		collection *mongo.Collection
	}
)

func (r *candidates) Add(ctx context.Context, candidate *entity.Candidate) error {
	if _, err := r.collection.InsertOne(ctx, candidate); err != nil {
		if isDuplicatedRowError(err) {
			return ErrDuplicatedRow
		}

		return err
	}

	return nil
}

func (r *candidates) GetByEmailAndCode(ctx context.Context, email, code string) (*entity.Candidate, error) {
	var result entity.Candidate
	if err := r.collection.FindOne(ctx, bson.M{"email": email, "code": code}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
