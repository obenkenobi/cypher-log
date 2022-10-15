package queryreq

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"strconv"
)

type ReqQueryReader interface {
	ReadString(key string, value *string) ReqQueryReader
	ReadStringOrDefault(key string, value *string, defaultValue string) ReqQueryReader
	ReadInt(key string, value *int) ReqQueryReader
	ReadIntOrDefault(key string, value *int, defaultValue int) ReqQueryReader
	ReadInt64(key string, value *int64) ReqQueryReader
	ReadInt64OrDefault(key string, value *int64, defaultValue int64) ReqQueryReader
	Complete() error
}

type GinCtxReqQueryReaderImpl struct {
	requiredQueryRuleErr apperrors.RuleError
	stringReadRuleErr    apperrors.RuleError
	intReadRuleErr       apperrors.RuleError
	int64ReadRuleErr     apperrors.RuleError
	ginContext           *gin.Context
	ruleErrors           []apperrors.RuleError
}

func (q GinCtxReqQueryReaderImpl) ReadString(key string, value *string) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		*value = queryVal
	} else {
		queryReader.ruleErrors = append(q.ruleErrors, q.requiredQueryRuleErr)
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

func (q GinCtxReqQueryReaderImpl) ReadInt(key string, value *int) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := strconv.Atoi(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(q.ruleErrors, q.intReadRuleErr)
		}
	} else {
		queryReader.ruleErrors = append(q.ruleErrors, q.requiredQueryRuleErr)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadIntOrDefault(key string, value *int, defaultValue int) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := strconv.Atoi(queryVal); err == nil {
			*value = intVal
		} else {
			queryReader.ruleErrors = append(q.ruleErrors, q.intReadRuleErr)
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
			queryReader.ruleErrors = append(q.ruleErrors, q.int64ReadRuleErr)
		}
	} else {
		queryReader.ruleErrors = append(q.ruleErrors, q.requiredQueryRuleErr)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) ReadInt64OrDefault(key string, value *int64, defaultValue int64) ReqQueryReader {
	queryReader := q
	if queryVal, ok := q.ginContext.GetQuery(key); ok {
		if intVal, err := utils.StrToInt64(queryVal); err == nil {
			*value = intVal
		} else {
			*value = defaultValue
		}
	} else {
		queryReader.ruleErrors = append(q.ruleErrors, q.requiredQueryRuleErr)
	}
	return queryReader
}

func (q GinCtxReqQueryReaderImpl) Complete() error {
	if len(q.ruleErrors) > 0 {
		return apperrors.NewBadReqErrorFromRuleErrors(q.ruleErrors...)
	}
	return nil
}

func NewGinCtxReqQueryReaderImpl(
	ginContext *gin.Context,
	requiredQueryRuleErr apperrors.RuleError,
	stringReadRuleErr apperrors.RuleError,
	intReadRuleErr apperrors.RuleError,
	int64ReadRuleErr apperrors.RuleError,
) GinCtxReqQueryReaderImpl {
	return GinCtxReqQueryReaderImpl{
		ginContext:           ginContext,
		requiredQueryRuleErr: requiredQueryRuleErr,
		stringReadRuleErr:    stringReadRuleErr,
		intReadRuleErr:       intReadRuleErr,
		int64ReadRuleErr:     int64ReadRuleErr,
	}
}
