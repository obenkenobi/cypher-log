type Direction = 'asc' | 'desc'

interface SortField {
  field: string,
  direction: Direction
}

interface PageRequest {
  page: bigint
  size: bigint
  sort: SortField[]
}

interface Page<T> {
  contents: T[],
  total: bigint
}