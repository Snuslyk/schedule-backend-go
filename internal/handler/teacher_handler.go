package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

type TeacherHandler struct {
	teacherRepo *postgres.TeacherRepository
}

func NewTeacherHandler(teacherRepo *postgres.TeacherRepository) *TeacherHandler {
	return &TeacherHandler{teacherRepo: teacherRepo}
}

// FindAll godoc
// @Summary      List all teachers
// @Tags         teacher
// @Produce      json
// @Success      200  {array}   domain.Teacher
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /api/teacher [get]
func (h *TeacherHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.teacherRepo.FindAll(r.Context())
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	JSON(w, http.StatusOK, teachers)
}

// FindOne godoc
// @Summary      Get teacher by ID
// @Tags         teacher
// @Produce      json
// @Param        id   path      int  true  "Teacher ID"
// @Success      200  {object}  domain.Teacher
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/teacher/{id} [get]
func (h *TeacherHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid teacher id")
		return
	}

	teacher, err := h.teacherRepo.FindByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if teacher == nil {
		NotFound(w, "Teacher not found")
		return
	}

	JSON(w, http.StatusOK, teacher)
}

// Create godoc
// @Summary      Create teacher
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateTeacherDTO  true  "Teacher payload"
// @Success      201    {object}  domain.Teacher
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/teacher [post]
func (h *TeacherHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateTeacherDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Name == "" {
		BadRequest(w, "Teacher name is required")
		return
	}

	t := &domain.Teacher{Name: dto.Name}
	if err := h.teacherRepo.Create(r.Context(), t); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, t)
}

// CreateMany godoc
// @Summary      Bulk create teachers
// @Tags         teacher
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateTeacherDTO  true  "List of teachers"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/teacher/bulk [post]
func (h *TeacherHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateTeacherDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.teacherRepo.CreateMany(r.Context(), dtos); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}

// Remove godoc
// @Summary      Delete teacher by ID
// @Tags         teacher
// @Param        id   path  int  true  "Teacher ID"
// @Success      204  "No Content"
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/teacher/{id} [delete]
func (h *TeacherHandler) Remove(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid teacher id")
		return
	}

	deleted, err := h.teacherRepo.DeleteByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, "Teacher not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
