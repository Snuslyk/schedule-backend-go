package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
)

type LessonHandler struct {
	lessonRepo *postgres.LessonRepository
}

func NewLessonHandler(lessonRepo *postgres.LessonRepository) *LessonHandler {
	return &LessonHandler{lessonRepo: lessonRepo}
}

// Create godoc
// @Summary      Create flat lesson
// @Description  Create a lesson with group_id and teacher_id foreign keys
// @Tags         lesson
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateLessonDTO  true  "Lesson payload"
// @Success      201    {object}  domain.Lesson
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/lesson [post]
func (h *LessonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateLessonDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.GroupID <= 0 || dto.TeacherID <= 0 || dto.SubjectID <= 0 {
		BadRequest(w, "groupId, teacherId, and subjectId are required")
		return
	}

	isAvailable := true
	if dto.IsAvailable != nil {
		isAvailable = *dto.IsAvailable
	}
	slotLength := dto.SlotLength
	if slotLength <= 0 {
		slotLength = 1
	}
	slotNumber := dto.SlotNumber
	if slotNumber <= 0 {
		slotNumber = 1
	}
	weekType := dto.WeekType
	if weekType == "" {
		weekType = domain.WeekTypeEven
	}

	l := &domain.Lesson{
		GroupID:     dto.GroupID,
		TeacherID:   dto.TeacherID,
		SubjectID:   dto.SubjectID,
		Classroom:   dto.Classroom,
		DayOfWeek:   dto.DayOfWeek,
		SlotNumber:  slotNumber,
		SlotLength:  slotLength,
		WeekType:    weekType,
		IsAvailable: isAvailable,
	}

	if err := h.lessonRepo.Create(r.Context(), l); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, l)
}

// CreateMany godoc
// @Summary      Bulk create lessons
// @Tags         lesson
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateLessonDTO  true  "List of lessons"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/lesson/bulk [post]
func (h *LessonHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateLessonDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.lessonRepo.CreateMany(r.Context(), dtos); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}

// DeleteByID godoc
// @Summary      Delete lesson by ID
// @Tags         lesson
// @Param        id   path  int  true  "Lesson ID"
// @Success      204  "No Content"
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/lesson/{id} [delete]
func (h *LessonHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid lesson id")
		return
	}

	deleted, err := h.lessonRepo.DeleteByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, "Lesson not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
