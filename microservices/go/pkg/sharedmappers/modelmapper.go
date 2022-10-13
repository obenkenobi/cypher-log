package sharedmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded"
)

func MapMongoModelToBaseId[M datasource.MongoModel](source M, dest *embedded.BaseId) {
	dest.Id = source.GetIdStr()
}

func MapMongoModelToBaseTimestamp[M datasource.MongoModel](source M, dest *embedded.BaseTimestamp) {
	dest.CreatedAt = source.GetCreatedAt().UnixMilli()
	dest.UpdatedAt = source.GetCreatedAt().UnixMilli()
}

func MapMongoModelToBaseCrudObject[M datasource.MongoModel](source M, dest *embedded.BaseCRUDObject) {
	MapMongoModelToBaseId(source, &dest.BaseId)
	MapMongoModelToBaseTimestamp(source, &dest.BaseTimestamp)

}
