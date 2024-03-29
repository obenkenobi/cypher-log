package mgmtools

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func ConvertSortDirection(dir pagination.Direction) int {
	if strings.EqualFold(string(dir), string(pagination.Descending)) {
		return -1
	}
	return 1
}

func SortFieldToMongoField(sortField string) string {
	switch sortField {
	case pagination.SortFieldCreatedAt:
		return "created_at"
	case pagination.SortFieldUpdatedAt:
		return "updated_at"
	default:
		return sortField
	}
}

func CreatePaginatedFindOpts(pageReq pagination.PageRequest) *options.FindOptions {
	return SetPaginatedFindOpts(options.Find(), pageReq)
}

func SetPaginatedFindOpts(findOpt *options.FindOptions, pageReq pagination.PageRequest) *options.FindOptions {
	findOpt = findOpt.SetSkip(pageReq.SkipCount()).SetLimit(pageReq.Size)
	for _, s := range pageReq.Sort {
		findOpt = findOpt.SetSort(bson.D{{SortFieldToMongoField(s.Field), ConvertSortDirection(s.Direction)}})
	}
	return findOpt
}

// HandleFindManyRes handles the result of a find many method of *mgm.Collection
// and transforms it into an observable.
func HandleFindManyRes[T any](ctx context.Context, cursor *mongo.Cursor, err error) ([]T, error) {
	var results []T
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	return results, err
}
