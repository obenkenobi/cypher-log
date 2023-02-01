/*
type BaseId struct {
	Id string `json:"id"`
}

type BaseRequiredId struct {
	Id string `json:"id" binding:"required"`
}

type BaseTimestamp struct {
	CreatedAt int64 `json:"createdAt"` // In unix timestamp in milliseconds
	UpdatedAt int64 `json:"updatedAt"` // In unix timestamp in milliseconds
}

type BaseCRUDObject struct {
	BaseId
	BaseTimestamp
}

* */

interface BaseId {
  id: string
}

interface BaseRequiredId {
  id: string
}
interface BaseTimestamp {
  createdAt: bigint,
  updatedAt: bigint
}

interface BaseCRUDObject extends BaseId, BaseTimestamp {}

