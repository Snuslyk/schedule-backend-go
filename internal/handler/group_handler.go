package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	groupRepo *postgres.GroupRepository
}

func NewGroupHandler(groupRepo *postgres.GroupRepository) *GroupHandler {
	return &GroupHandler{groupRepo: groupRepo}
}

// GetByID godoc
// @Summary      Get group by ID
// @Tags         group
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  domain.Group
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/group/id/{id} [get]
func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid group id")
		return
	}

	group, err := h.groupRepo.FindByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if group == nil {
		NotFound(w, "Group not found")
		return
	}

	JSON(w, http.StatusOK, group)
}

// GetByName godoc
// @Summary      Get group by name
// @Tags         group
// @Produce      json
// @Param        name  path      string  true  "Group Name"
// @Success      200   {object}  domain.Group
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Router       /api/group/name/{name} [get]
func (h *GroupHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if unescaped, err := url.PathUnescape(name); err == nil && unescaped != "" {
		name = unescaped
	}
	if name == "" {
		BadRequest(w, "Group name is required")
		return
	}

	group, err := h.groupRepo.FindByName(r.Context(), name)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if group == nil {
		NotFound(w, "Group not found")
		return
	}

	JSON(w, http.StatusOK, group)
}

// Create godoc
// @Summary      Create group
// @Tags         group
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateGroupDTO  true  "Group payload"
// @Success      201    {object}  domain.Group
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      409    {object}  domain.ErrorResponse
// @Router       /api/group [post]
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Name == "" {
		BadRequest(w, "Group name is required")
		return
	}

	g := &domain.Group{Name: dto.Name}
	if err := h.groupRepo.Create(r.Context(), g); err != nil {
		Conflict(w, "Group with this name already exists")
		return
	}

	JSON(w, http.StatusCreated, g)
}

// CreateMany godoc
// @Summary      Bulk create groups
// @Tags         group
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateGroupDTO  true  "List of groups"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      409    {object}  domain.ErrorResponse
// @Router       /api/group/bulk [post]
func (h *GroupHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.groupRepo.CreateMany(r.Context(), dtos); err != nil {
		Conflict(w, "One or more groups already exist")
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}

// DeleteByID godoc
// @Summary      Delete group by ID
// @Tags         group
// @Param        id   path  int  true  "Group ID"
// @Success      204  "No Content"
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/group/id/{id} [delete]
func (h *GroupHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid group id")
		return
	}

	deleted, err := h.groupRepo.DeleteByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, "Group not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteByName godoc
// @Summary      Delete group by name
// @Tags         group
// @Param        name  path  string  true  "Group Name"
// @Success      204   "No Content"
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Router       /api/group/name/{name} [delete]
func (h *GroupHandler) DeleteByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if unescaped, err := url.PathUnescape(name); err == nil && unescaped != "" {
		name = unescaped
	}
	if name == "" {
		BadRequest(w, "Group name is required")
		return
	}

	deleted, err := h.groupRepo.DeleteByName(r.Context(), name)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, "Group not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
