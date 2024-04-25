package repository

//go:generate mockery --name=Tx --inpackage --output=. --filename=tx_mock.go --structname=TxMock

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewTx(client *mongo.Client) Tx {
	return &tx{client: client}
}

type (
	Tx interface {
		Do(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) error
	}

	tx struct {
		client *mongo.Client
	}
)

func (t *tx) Do(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	var (
		err  error
		sess mongo.Session
	)
	if sess, err = t.client.StartSession(); err != nil {
		return fmt.Errorf("can't start session: %w", err)
	}
	defer sess.EndSession(ctx)

	if _, err = sess.WithTransaction(ctx, callback); err != nil {
		return fmt.Errorf("can't execute transaction: %w", err)
	}

	return nil
}
