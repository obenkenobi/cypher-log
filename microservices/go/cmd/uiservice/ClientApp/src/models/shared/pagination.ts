type Direction = 'asc' | 'desc'

interface SortField {
  field: string,
  direction: Direction
}

interface PageRequest {
  page: number
  size: number
  sort: SortField[]
}

interface Page<T> {
  contents: T[],
  total: number
}