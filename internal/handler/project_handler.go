package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"todo-api/internal/service"
)

type ProjectHandler struct {
	projects *service.ProjectService
}

func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

type projectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	project, err := h.projects.Create(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	result, err := h.projects.List(c.Request.Context(), page, pageSize)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProjectHandler) Stats(c *gin.Context) {
	stats, err := h.projects.Stats(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	project, err := h.projects.GetByID(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	project, err := h.projects.Update(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, ok := parsePositiveID(c, "id")
	if !ok {
		return
	}
	if err := h.projects.Delete(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
