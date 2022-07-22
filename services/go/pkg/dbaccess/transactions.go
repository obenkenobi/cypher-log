package dbaccess

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Session interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

type TransactionRunner interface {
	ExecTransaction(runner func(Session, context.Context) error) error
}

type TransactionRunnerMongo struct {
}

func NewTransactionRunnerMongo(dbClient DBMongoClient) TransactionRunner {
	return &TransactionRunnerMongo{}
}

func (u TransactionRunnerMongo) ExecTransaction(transactionFunc func(Session, context.Context) error) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		return transactionFunc(session, sc)
	})
}
