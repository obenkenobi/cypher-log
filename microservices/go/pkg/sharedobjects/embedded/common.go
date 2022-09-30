package embedded

type BaseId struct {
	Id string `json:"id"`
}

type BaseTimestamp struct {
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}
