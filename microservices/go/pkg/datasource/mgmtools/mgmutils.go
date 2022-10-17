package mgmtools

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
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

func CreatePaginatedFindOpts(pageReq pagination.PageRequest) *options.FindOptions {
	return SetPaginatedFindOpts(options.Find(), pageReq)
}

func SetPaginatedFindOpts(findOpt *options.FindOptions, pageReq pagination.PageRequest) *options.FindOptions {
	findOpt = findOpt.SetSkip(pageReq.SkipCount()).SetLimit(pageReq.Size)
	for _, s := range pageReq.Sort {
		findOpt = findOpt.SetSort(bson.D{{s.Field, ConvertSortDirection(s.Direction)}})
	}
	return findOpt
}

// HandleFindManyRes handles the result of a find many method of *mgm.Collection
// and transforms it into an observable.
func HandleFindManyRes[T any](ctx context.Context, cursor *mongo.Cursor, err error) stream.Observable[T] {
	var results []T
	if err != nil {
		return stream.Error[models.Note](err)
	}
	if err = cursor.All(ctx, &results); err != nil {
		return stream.Error[models.Note](err)
	}
	return stream.FromSlice(results)
}
