package repositories

import "github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"

type UserRepository interface {
}

type UserRepositoryImpl struct {
}

func NewUserRepositoryImpl(mongoClient dbaccess.DBMongoClient) *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}
