package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type WeekTemplateHandler struct {
	wtService *service.WeekTemplateService
}

func NewWeekTemplateHandler(wtService *service.WeekTemplateService) *WeekTemplateHandler {
	return &WeekTemplateHandler{wtService: wtService}
}

// Create godoc
// @Summary      Create a week template for a group
// @Tags         week-template
// @Accept       json
// @Produce      json
// @Param        id     path      int                           true  "Group ID"
// @Param        input  body      domain.CreateWeekTemplateDTO  true  "Week template payload"
// @Success      201    {object}  domain.WeekTemplate
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      404    {object}  domain.ErrorResponse
// @Router       /api/week-template/{id} [post]
func (h *WeekTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid group id")
		return
	}

	var dto domain.CreateWeekTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	wt, err := h.wtService.Create(r.Context(), groupID, dto)
	if err != nil {
		if strings.Contains(err.Error(), "no group with id") {
			BadRequest(w, err.Error())
			return
		}
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, wt)
}

// CreateMany godoc
// @Summary      Create multiple week templates for a group
// @Tags         week-template
// @Accept       json
// @Produce      json
// @Param        id     path      int                             true  "Group ID"
// @Param        input  body      []domain.CreateWeekTemplateDTO  true  "List of week templates"
// @Success      201    {array}   domain.WeekTemplate
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/week-template/{id}/bulk [post]
func (h *WeekTemplateHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid group id")
		return
	}

	var dtos []domain.CreateWeekTemplateDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	results, err := h.wtService.CreateMany(r.Context(), groupID, dtos)
	if err != nil {
		if strings.Contains(err.Error(), "no group with id") {
			BadRequest(w, err.Error())
			return
		}
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, results)
}

// DeleteByID godoc
// @Summary      Delete a week template by ID
// @Tags         week-template
// @Param        id   path      int  true  "Week Template ID"
// @Success      204  "No Content"
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/week-template/{id} [delete]
func (h *WeekTemplateHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid week template id")
		return
	}

	deleted, err := h.wtService.DeleteByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, fmt.Sprintf("Week Template with id %d not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
