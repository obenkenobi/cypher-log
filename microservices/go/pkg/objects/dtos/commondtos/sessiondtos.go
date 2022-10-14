package commondtos

type UserKeySessionDto struct {
	ProxyKid      string `json:"proxyKid" binding:"required"`
	Token         string `json:"token" binding:"required"`
	UserId        string `json:"userId" binding:"required"`
	KeyVersion    int64  `json:"keyVersion"`
	StartTime     int64  `json:"startTime"`     // In unix timestamp in milliseconds
	DurationMilli int64  `json:"durationMilli"` // In milliseconds
}

type UserKeySessionPayloadDto[T any] struct {
	Session UserKeySessionDto `json:"session"`
	Payload T                 `json:"payload"`
}
