package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type SubjectHandler struct {
	subjectRepo *postgres.SubjectRepository
}

func NewSubjectHandler(subjectRepo *postgres.SubjectRepository) *SubjectHandler {
	return &SubjectHandler{subjectRepo: subjectRepo}
}

// FindAll godoc
// @Summary      List all subjects
// @Tags         subject
// @Produce      json
// @Success      200  {array}   domain.Subject
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /api/subject [get]
func (h *SubjectHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.subjectRepo.FindAll(r.Context())
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	JSON(w, http.StatusOK, subjects)
}

// Create godoc
// @Summary      Create subject
// @Tags         subject
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateSubjectDTO  true  "Subject payload"
// @Success      201    {object}  domain.Subject
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      409    {object}  domain.ErrorResponse
// @Router       /api/subject [post]
func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateSubjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Name == "" {
		BadRequest(w, "Subject name is required")
		return
	}

	s := &domain.Subject{Name: dto.Name}
	if err := h.subjectRepo.Create(r.Context(), s); err != nil {
		Conflict(w, "Subject with this name already exists")
		return
	}

	JSON(w, http.StatusCreated, s)
}

// CreateMany godoc
// @Summary      Bulk create subjects
// @Tags         subject
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateSubjectDTO  true  "List of subjects"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      409    {object}  domain.ErrorResponse
// @Router       /api/subject/bulk [post]
func (h *SubjectHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateSubjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.subjectRepo.CreateMany(r.Context(), dtos); err != nil {
		Conflict(w, "One or more subjects already exist")
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}
