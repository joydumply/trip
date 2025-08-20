package service

import (
	"context"
	"trip-planner/internal/model"
	"trip-planner/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrip(ctx context.Context, trip *model.Trip) error {
	return s.repo.CreateTrip(ctx, trip)
}

// old but gold

// func (s *Service) GetTrips(ctx context.Context) ([]model.Trip, error) {
// 	return s.repo.GetTrips(ctx)
// }

func (s *Service) GetTrips(ctx context.Context, title string, leaderID *int64, limit, offset int, sortBy, order string) ([]model.Trip, error) {
	return s.repo.GetTrips(ctx, title, leaderID, limit, offset, sortBy, order)
}

func (s *Service) GetTripByID(ctx context.Context, id int64) (*model.Trip, error) {
	return s.repo.GetTripByID(ctx, id)
}

func (s *Service) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	return s.repo.UpdateTrip(ctx, trip)
}

func (s *Service) DeleteTrip(ctx context.Context, id int64) error {

	return s.repo.DeleteTrip(ctx, id)
}
