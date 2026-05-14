package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/cache"
	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/Snuslyk/schedule-backend-go/internal/utils"
)

type ReplaceService struct {
	replaceRepo *postgres.ReplaceRepository
	groupRepo   *postgres.GroupRepository
	cache       *cache.ScheduleCache
}

func NewReplaceService(
	replaceRepo *postgres.ReplaceRepository,
	groupRepo *postgres.GroupRepository,
	cache *cache.ScheduleCache,
) *ReplaceService {
	return &ReplaceService{
		replaceRepo: replaceRepo,
		groupRepo:   groupRepo,
		cache:       cache,
	}
}

func (s *ReplaceService) Create(ctx context.Context, groupName string, dto domain.CreateReplaceDTO, thisWeek time.Time) (*domain.Replace, error) {
	group, err := s.groupRepo.FindByName(ctx, groupName)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group with name %s not found", groupName)
	}

	rep := &domain.Replace{
		Date:        dto.Date,
		SlotNumber:  dto.SlotNumber,
		GroupID:     group.ID,
		TeacherID:   dto.TeacherID,
		SubjectID:   dto.SubjectID,
		Classroom:   dto.Classroom,
		IsAvailable: dto.IsAvailable,
	}

	if err := s.replaceRepo.Create(ctx, rep); err != nil {
		return nil, err
	}

	now := time.Now()
	if !thisWeek.IsZero() {
		now = thisWeek
	}

	if utils.IsSameWeek(now, dto.Date) {
		s.cache.InvalidateGroup(ctx, group.ID, group.Name)
		s.cache.InvalidateTeacher(ctx, dto.TeacherID)
	}

	return rep, nil
}

func (s *ReplaceService) CreateMany(ctx context.Context, groupName string, dtos []domain.CreateReplaceDTO, thisWeek time.Time) error {
	group, err := s.groupRepo.FindByName(ctx, groupName)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group with name %s not found", groupName)
	}

	reps := make([]domain.Replace, len(dtos))
	now := time.Now()
	if !thisWeek.IsZero() {
		now = thisWeek
	}

	hasThisWeek := false
	teacherIDs := make(map[int]bool)

	for i, dto := range dtos {
		reps[i] = domain.Replace{
			Date:        dto.Date,
			SlotNumber:  dto.SlotNumber,
			GroupID:     group.ID,
			TeacherID:   dto.TeacherID,
			SubjectID:   dto.SubjectID,
			Classroom:   dto.Classroom,
			IsAvailable: dto.IsAvailable,
		}
		teacherIDs[dto.TeacherID] = true
		if utils.IsSameWeek(now, dto.Date) {
			hasThisWeek = true
		}
	}

	if err := s.replaceRepo.CreateMany(ctx, reps); err != nil {
		return err
	}

	if hasThisWeek {
		s.cache.InvalidateGroup(ctx, group.ID, group.Name)
		for tid := range teacherIDs {
			s.cache.InvalidateTeacher(ctx, tid)
		}
	}

	return nil
}
