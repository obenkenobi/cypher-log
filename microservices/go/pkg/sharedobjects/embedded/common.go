package embedded

type BaseId struct {
	Id string `json:"id"`
}

type BaseTimestamp struct {
	CreatedAt int64 `json:"createdAt"` // In unix timestamp in milliseconds
	UpdatedAt int64 `json:"updatedAt"` // In unix timestamp in milliseconds
}

type BaseCRUDObject struct {
	BaseId
	BaseTimestamp
}
