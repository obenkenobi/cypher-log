package dshandlers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type DataSourceHandler interface {
	// IsNotFoundError checks if an error is created by an underlying object database
	// mapper is due to a requested entity not being found.
	IsNotFoundError(err error) bool
}

func HandleSingleFind[TQueryResult any](
	dbHandler DataSourceHandler,
	supplier func() (TQueryResult, error),
) (option.Maybe[TQueryResult], error) {
	if result, err := supplier(); err != nil {
		if dbHandler.IsNotFoundError(err) {
			return option.None[TQueryResult](), nil
		}
		return nil, err
	} else {
		return option.Perhaps(result), nil
	}
}
