package apibase

import (
	"fmt"

	"github.com/kataras/iris/context"
)

type Pagination struct {
	Start    int
	Limit    int
	SortBy   string
	SortDesc bool
}

func (p *Pagination) Dir() string {
	if p.SortDesc {
		return "DESC"
	} else {
		return "ASC"
	}
}

func (p *Pagination) AsQuery() string {
	s := ""
	if p.SortBy != "" {
		s += "ORDER BY " + p.SortBy + " " + p.Dir()
	}
	if p.Limit > 0 {
		if p.Start > 0 {
			s += fmt.Sprintf(" LIMIT %d,%d", p.Start, p.Limit)
		} else {
			s += fmt.Sprintf(" LIMIT %d", p.Limit)
		}
	} else if p.Start > 0 {
		s += fmt.Sprintf(" LIMIT %d,10", p.Start)
	}
	return s
}

func ExtendPaginationQueryParamMap(paramsMap ApiParams) ApiParams {
	if paramsMap == nil {
		paramsMap = ApiParams{}
	}
	testAndSet := func(key string, param ApiParam) {
		if _, existed := paramsMap[key]; !existed {
			paramsMap[key] = param
		}
	}
	testAndSet("start", ApiParam{"int", "起始页", "0", false})
	testAndSet("limit", ApiParam{"int", "分页大小", "10", false})
	testAndSet("sort-by", ApiParam{"string", "排序依据", "", false})
	testAndSet("sort-dir", ApiParam{"string", "排序方向", "", false})
	return paramsMap
}

func GetPaginationFromQueryParameters(paramErrs *ApiParameterErrors, ctx context.Context) *Pagination {
	start, err := ctx.URLParamInt("start")
	if err != nil {
		if s := ctx.URLParam("start"); s == "" {
			start = 0
		} else {
			if paramErrs != nil {
				paramErrs.AppendError("start", "invalid value: %s", s)
			}
		}
	}
	limit, err := ctx.URLParamInt("limit")
	if err != nil {
		if s := ctx.URLParam("limit"); s == "" {
			limit = 10
		} else {
			if paramErrs != nil {
				paramErrs.AppendError("limit", "invalid value %s", s)
			}
		}
	}
	sortBy := ctx.URLParam("sort-by")
	sortDir := ctx.URLParam("sort-dir")
	sortDesc := false
	switch sortDir {
	case "1", "asc", "ascend", "true", "":
	case "0", "desc", "descend", "false":
		sortDesc = true
	default:
		if paramErrs != nil {
			paramErrs.AppendError("sort-dir", "invalid value: %s", sortDir)
		}
	}

	return &Pagination{start, limit, sortBy, sortDesc}
}
