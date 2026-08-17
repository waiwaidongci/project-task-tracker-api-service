package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todo-api/internal/repository"
	"todo-api/internal/service"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{Error: ErrorDetail{Code: code, Message: message}})
}

func respondServiceError(c *gin.Context, err error) {
	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		respondError(c, http.StatusBadRequest, "validation_error", validationErr.Error())
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "not_found", "resource not found")
		return
	}
	respondError(c, http.StatusInternalServerError, "internal_error", err.Error())
}

func parsePositiveID(c *gin.Context, key string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "validation_error", key+" must be a positive integer")
		return 0, false
	}
	return id, true
}

func parsePagination(c *gin.Context) (int, int, bool) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 20)
	if page < 1 || pageSize < 1 || pageSize > 100 {
		respondError(c, http.StatusBadRequest, "validation_error", "page must be >= 1 and page_size must be between 1 and 100")
		return 0, 0, false
	}
	return page, pageSize, true
}

func queryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
