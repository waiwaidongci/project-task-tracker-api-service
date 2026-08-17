package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todo-api/internal/service"
)

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

type taskRequest struct {
	ProjectID int64    `json:"project_id"`
	Title     string   `json:"title"`
	Priority  string   `json:"priority"`
	Status    string   `json:"status"`
	DueDate   string   `json:"due_date"`
	Tags      []string `json:"tags"`
}

type statusRequest struct {
	Status string `json:"status"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	task, err := h.tasks.Create(c.Request.Context(), service.TaskInput{
		ProjectID: req.ProjectID,
		Title:     req.Title,
		Priority:  req.Priority,
		Status:    req.Status,
		DueDate:   req.DueDate,
		Tags:      req.Tags,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	filter := service.TaskListFilter{
		ProjectID: int64(queryInt(c, "project_id", 0)),
		Status:    c.Query("status"),
		Priority:  c.Query("priority"),
		DueToday:  c.Query("due_today") == "true",
	}
	result, err := h.tasks.List(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) Today(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	result, err := h.tasks.Today(c.Request.Context(), page, pageSize)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	task, err := h.tasks.GetByID(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	task, err := h.tasks.Update(c.Request.Context(), id, service.TaskInput{
		ProjectID: req.ProjectID,
		Title:     req.Title,
		Priority:  req.Priority,
		Status:    req.Status,
		DueDate:   req.DueDate,
		Tags:      req.Tags,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ChangeStatus(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	task, err := h.tasks.ChangeStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := h.tasks.Delete(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseBoolQuery(c *gin.Context, key string) bool {
	value := c.Query(key)
	if value == "" {
		return false
	}
	parsed, err := strconv.ParseBool(value)
	return err == nil && parsed
}
