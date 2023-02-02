interface BaseId {
  id: string
}

interface BaseRequiredId {
  id: string
}
interface BaseTimestamp {
  createdAt: number,
  updatedAt: number
}

interface BaseCRUDObject extends BaseId, BaseTimestamp {}

