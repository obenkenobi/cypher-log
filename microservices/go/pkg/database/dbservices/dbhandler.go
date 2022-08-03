package dbservices

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

type DBHandler interface {
	// IsNotFoundError checks if an error is created by an underlying object database
	// mapper is due to a requested entity not being found.
	IsNotFoundError(err error) bool
}

// ObserveOptionalSingleQueryAsync
//creates a single out of a supplier function that queries a single value. The
//supplier function is run on a separate goroutine. *Make sure your supplier
//function is not going to be thread safe or not cause race conditions on the
//data accessed.
func ObserveOptionalSingleQueryAsync[TQueryResult any](
	dbHandler DBHandler,
	supplier func() (TQueryResult, error),
) single.Single[option.Maybe[TQueryResult]] {
	return single.FromSupplierAsync(func() (option.Maybe[TQueryResult], error) {
		return runOptionalSingleQuery(dbHandler, supplier)
	})
}

func ObserveOptionalSingleQuery[TQueryResult any](
	dbHandler DBHandler,
	supplier func() (TQueryResult, error),
) single.Single[option.Maybe[TQueryResult]] {
	return single.FromSupplier(func() (option.Maybe[TQueryResult], error) {
		return runOptionalSingleQuery(dbHandler, supplier)
	})
}

func runOptionalSingleQuery[TQueryResult any](
	dbHandler DBHandler,
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
