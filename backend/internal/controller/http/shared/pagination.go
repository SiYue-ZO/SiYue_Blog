package shared

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

type PaginationParams struct {
	Page     int
	PageSize int
}

func ParsePagination(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = defaultPage
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}
