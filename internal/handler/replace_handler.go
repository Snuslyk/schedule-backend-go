package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/Snuslyk/schedule-backend-go/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ReplaceHandler struct {
	replaceService *service.ReplaceService
}

func NewReplaceHandler(replaceService *service.ReplaceService) *ReplaceHandler {
	return &ReplaceHandler{replaceService: replaceService}
}

// Create godoc
// @Summary      Create replace for group
// @Description  Create replace and invalidate cache in Redis for group and teacher
// @Tags         replace
// @Accept       json
// @Produce      json
// @Param        group     path      string                   true   "Group name"
// @Param        thisWeek  query     string                   false  "Date in this week YYYY-MM-DD"
// @Param        input     body      domain.CreateReplaceDTO  true   "Replace payload"
// @Success      201       {object}  domain.Replace
// @Failure      400       {object}  domain.ErrorResponse
// @Router       /api/replace/{group} [post]
func (h *ReplaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "group")
	if unescaped, err := url.PathUnescape(groupName); err == nil && unescaped != "" {
		groupName = unescaped
	}
	if groupName == "" {
		BadRequest(w, "Group name is required")
		return
	}

	thisWeekStr := r.URL.Query().Get("thisWeek")
	var thisWeek time.Time
	if thisWeekStr != "" {
		thisWeek, _ = utils.ParseDate(thisWeekStr)
	}

	var dto domain.CreateReplaceDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	rep, err := h.replaceService.Create(r.Context(), groupName, dto, thisWeek)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, rep)
}

// CreateMany godoc
// @Summary      Bulk create replaces for group
// @Tags         replace
// @Accept       json
// @Produce      json
// @Param        group     path      string                     true   "Group name"
// @Param        thisWeek  query     string                     false  "Date in this week YYYY-MM-DD"
// @Param        input     body      []domain.CreateReplaceDTO  true   "List of replaces"
// @Success      201       {object}  map[string]int
// @Failure      400       {object}  domain.ErrorResponse
// @Router       /api/replace/{group}/bulk [post]
func (h *ReplaceHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "group")
	if unescaped, err := url.PathUnescape(groupName); err == nil && unescaped != "" {
		groupName = unescaped
	}
	if groupName == "" {
		BadRequest(w, "Group name is required")
		return
	}

	thisWeekStr := r.URL.Query().Get("thisWeek")
	var thisWeek time.Time
	if thisWeekStr != "" {
		thisWeek, _ = utils.ParseDate(thisWeekStr)
	}

	var dtos []domain.CreateReplaceDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.replaceService.CreateMany(r.Context(), groupName, dtos, thisWeek); err != nil {
		BadRequest(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}
