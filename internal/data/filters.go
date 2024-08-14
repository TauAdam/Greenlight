package data

import (
	"github.com/TauAdam/Greenlight/internal/validator"
	"math"
	"strings"
)

type Filters struct {
	Page        int
	PageSize    int
	Sort        string
	SortOptions []string
}

// sortColumn returns the column name to sort the database query results.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortOptions {
		if f.Sort != safeValue {
			panic("unsafe sort parameter: " + f.Sort)
		}
	}
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000, "page", "must be a maximum of 10 thousands")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortOptions...), "sort", "invalid sort value")
}

type Metadata struct {
	FirstPage    int `json:"first_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	CurrentPage  int `json:"current_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
}

// calculateMetadata calculates the metadata for a paginated response.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		FirstPage:    1,
		PageSize:     pageSize,
		CurrentPage:  page,
		TotalRecords: totalRecords,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
	}
}
