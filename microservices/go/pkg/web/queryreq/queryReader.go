package queryreq

import (
	"github.com/akrennmair/slice"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"strconv"
	"strings"
)

type ReqQueryReader interface {
	ReadString(key string, value *string) ReqQueryReader
	ReadStringOrDefault(key string, value *string, defaultValue string) ReqQueryReader
	ReadBoolean(key string, value *bool) ReqQueryReader
	ReadBoolOrDefault(key string, value *bool, defaultValue bool) ReqQueryReader
	ReadInt(key string, value *int) ReqQueryReader
	ReadIntOrDefault(key string, value *int, defaultValue int) ReqQueryReader
	ReadInt64(key string, value *int64) ReqQueryReader
	ReadInt64OrDefault(key string, value *int64, defaultValue int64) ReqQueryReader
	ReadSort(key string, value *[]pagination.SortField) ReqQueryReader
	ReadSortOrDefault(key string, value *[]pagination.SortField, defaultValue []pagination.SortField) ReqQueryReader
	ReadPageRequest(pageKey string, sizeKey string, sortKey string, request *pagination.PageRequest) ReqQueryReader
	ReadPageRequestOrDefault(
		pageKey string,
		sizeKey string,
		sortKey string,
		request *pagination.PageRequest,
		defaultValue pagination.PageRequest,
	) ReqQueryReader
	Complete() error
}

type GinCtxReqQueryReaderImpl struct {
	errorService sharedservices.ErrorService
	ginContext   *gin.Context
	ruleErrors   []apperrors.RuleError
}

func (q GinCtxReqQueryReaderImpl) ReadString(key string, value *string) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		*value = queryVal
	} else {
		queryReader.ruleErrors = append(
			q.ruleErrors,
			q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryRequired, key),
		)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadStringOrDefault(key string, value *string, defaultValue string) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		*value = queryVal
	} else {
		*value = defaultValue
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadBoolean(key string, value *bool) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		queryValLower := strings.ToLower(queryVal)
		switch queryValLower {
		case "true":
			*value = true
		case "false":
			*value = false
		default:
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryBoolParseFail, key),
			)
		}
	} else {
		queryReader.ruleErrors = append(
			q.ruleErrors,
			q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryRequired, key),
		)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadBoolOrDefault(key string, value *bool, defaultValue bool) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		queryValLower := strings.ToLower(queryVal)
		switch queryValLower {
		case "true":
			*value = true
		case "false":
			*value = false
		default:
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryBoolParseFail, key),
			)
		}
	} else {
		*value = defaultValue
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadInt(key string, value *int) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := strconv.Atoi(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryIntParseFail, key),
			)
		}
	} else {
		queryReader.ruleErrors = append(
			q.ruleErrors,
			q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryRequired, key),
		)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadIntOrDefault(key string, value *int, defaultValue int) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := strconv.Atoi(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryIntParseFail, key),
			)
		}
	} else {
		*value = defaultValue
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadInt64(key string, value *int64) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := utils.StrToInt64(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryIntParseFail, key),
			)
		}
	} else {
		queryReader.ruleErrors = append(
			q.ruleErrors,
			q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryRequired, key),
		)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadInt64OrDefault(key string, value *int64, defaultValue int64) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := utils.StrToInt64(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(
				q.ruleErrors,
				q.errorService.RuleErrorFromCode(apperrors.ErrCodeReqQueryIntParseFail, key),
			)
		}
	} else {
		*value = defaultValue
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadSort(key string, value *[]pagination.SortField) ReqQueryReader {
	queryReader := q
	queryValues, _ := q.ginContext.GetQueryArray(key)
	sortFields := q.parseSortQueryStrings(queryValues)
	*value = sortFields
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadSortOrDefault(
	key string,
	value *[]pagination.SortField,
	defaultValue []pagination.SortField,
) ReqQueryReader {
	queryReader := q
	queryValues, ok := q.ginContext.GetQueryArray(key)
	if ok {
		sortFields := q.parseSortQueryStrings(queryValues)
		*value = sortFields
	} else {
		*value = defaultValue
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) parseSortQueryStrings(values []string) []pagination.SortField {
	mappedSortFields := slice.Map(values, func(v string) pagination.SortField {
		splitValues := strings.SplitN(v, ",", 2)
		if len(splitValues) == 0 {
			return pagination.NewSortField("", pagination.Ascending)
		}
		field := splitValues[0]
		var direction pagination.Direction
		if len(splitValues) > 1 && strings.EqualFold(splitValues[1], string(pagination.Descending)) {
			direction = pagination.Descending
		} else {
			direction = pagination.Ascending
		}
		return pagination.NewSortField(field, direction)
	})
	return slice.Filter(mappedSortFields, func(sf pagination.SortField) bool {
		return utils.StringIsBlank(sf.Field)
	})
}

func (q GinCtxReqQueryReaderImpl) ReadPageRequest(
	pageKey string,
	sizeKey string,
	sortKey string,
	request *pagination.PageRequest,
) ReqQueryReader {
	return q.ReadInt64(pageKey, &(request.Page)).
		ReadInt64(sizeKey, &(request.Size)).
		ReadSort(sortKey, &(request.Sort))
}

func (q GinCtxReqQueryReaderImpl) ReadPageRequestOrDefault(
	pageKey string,
	sizeKey string,
	sortKey string,
	request *pagination.PageRequest,
	defaultReq pagination.PageRequest,
) ReqQueryReader {
	return q.ReadInt64OrDefault(pageKey, &(request.Page), defaultReq.Page).
		ReadInt64OrDefault(sizeKey, &(request.Size), defaultReq.Size).
		ReadSortOrDefault(sortKey, &(request.Sort), defaultReq.Sort)
}

func (q GinCtxReqQueryReaderImpl) Complete() error {
	if len(q.ruleErrors) > 0 {
		return apperrors.NewBadReqErrorFromRuleErrors(q.ruleErrors...)
	}
	return nil
}

func NewGinCtxReqQueryReaderImpl(
	ginContext *gin.Context,
	errorService sharedservices.ErrorService,
) GinCtxReqQueryReaderImpl {
	return GinCtxReqQueryReaderImpl{
		ginContext:   ginContext,
		errorService: errorService,
	}
}
