package util

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {
	if page <= 0 {
		page = 1
	}

	defaultLimit := GetEnvInt("DEFAULT_PAGE_SIZE", 10)
	maxLimit := GetEnvInt("MAX_PAGE_SIZE", 100)

	if defaultLimit <= 0 {
		defaultLimit = 10
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}

	if limit <= 0 {
		limit = int32(defaultLimit)
	}

	if limit > int32(maxLimit) {
		limit = int32(maxLimit)
	}

	if limit <= 0 {
		limit = 10
	}

	totalPages := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}
