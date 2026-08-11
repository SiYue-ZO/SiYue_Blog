package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedData struct {
	List       interface{} `json:"list"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, bizCode int, message string) {
	c.JSON(httpStatus, Envelope{
		Code:    bizCode,
		Message: message,
		Data:    nil,
	})
}

func Paginated(c *gin.Context, list interface{}, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, Envelope{
		Code:    0,
		Message: "success",
		Data: PaginatedData{
			List:       list,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
