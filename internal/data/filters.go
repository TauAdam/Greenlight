package data

import (
	"github.com/TauAdam/Greenlight/internal/validator"
	"strings"
)

type Filters struct {
	Page        int
	PageSize    int
	Sort        string
	SortOptions []string
}

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

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000, "page", "must be a maximum of 10 thousands")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortOptions...), "sort", "invalid sort value")
}
