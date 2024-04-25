package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDuplicatedRow = errors.New("duplicated row")

func isDuplicatedRowError(err error) bool {
	var writeException mongo.WriteException
	if errors.As(err, &writeException) {
		for _, writeErr := range writeException.WriteErrors {
			if writeErr.Code == 11000 {
				return true
			}
		}
	}
	return false
}
