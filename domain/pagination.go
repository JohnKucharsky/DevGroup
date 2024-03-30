package domain

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginationParams struct {
	Limit  string `json:"limit" validate:"omitempty,numeric"`
	Offset string `json:"offset" validate:"omitempty,numeric"`
}

type ParsedPaginationParams struct {
	Limit  int
	Offset *int
}
