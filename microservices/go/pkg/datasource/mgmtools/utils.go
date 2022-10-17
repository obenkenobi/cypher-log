package mgmtools

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"strings"
)

func ConvertSortDirection(dir pagination.Direction) int {
	if strings.EqualFold(string(dir), string(pagination.Descending)) {
		return -1
	}
	return 1
}
