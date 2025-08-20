package service

import (
	"context"
	"fmt"
	"testing"
	"trip-planner/internal/model"
)

// --- Mock репозиторий ---
type MockRepo struct {
	trips map[int64]*model.Trip
}

func NewMockRepo() *MockRepo {
	return &MockRepo{trips: make(map[int64]*model.Trip)}
}

func (m *MockRepo) CreateTrip(ctx context.Context, trip *model.Trip) error {
	trip.ID = int64(len(m.trips) + 1)
	m.trips[trip.ID] = trip
	return nil
}

func (m *MockRepo) GetTrips(ctx context.Context) ([]model.Trip, error) {
	list := []model.Trip{}
	for _, t := range m.trips {
		list = append(list, *t)
	}
	return list, nil
}

func (m *MockRepo) GetTripByID(ctx context.Context, id int64) (*model.Trip, error) {
	if trip, ok := m.trips[id]; ok {
		return trip, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *MockRepo) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	if _, ok := m.trips[trip.ID]; !ok {
		return fmt.Errorf("trip with id %d not found", trip.ID)
	}
	m.trips[trip.ID] = trip
	return nil
}

func (m *MockRepo) DeleteTrip(ctx context.Context, id int64) error {
	if _, ok := m.trips[id]; !ok {
		return fmt.Errorf("trip with id %d not found", id)
	}
	delete(m.trips, id)
	return nil
}

func (m *MockRepo) Close() {}

// --- Тесты сервиса ---
func TestServiceCRUD(t *testing.T) {
	repo := NewMockRepo()
	srv := NewService(repo)

	ctx := context.Background()

	// --- Create ---
	trip := &model.Trip{
		LeaderID:    1,
		Title:       "Trip 1",
		Description: "Desc 1",
	}
	err := srv.CreateTrip(ctx, trip)
	if err != nil {
		t.Fatalf("CreateTrip failed: %v", err)
	}
	if trip.ID == 0 {
		t.Fatal("expected trip.ID to be assigned")
	}

	// --- GetTrips ---
	trips, err := srv.GetTrips(ctx)
	if err != nil {
		t.Fatalf("GetTrips failed: %v", err)
	}
	if len(trips) != 1 {
		t.Fatalf("expected 1 trip, got %d", len(trips))
	}

	// --- GetTripByID ---
	got, err := srv.GetTripByID(ctx, trip.ID)
	if err != nil {
		t.Fatalf("GetTripByID failed: %v", err)
	}
	if got.Title != trip.Title {
		t.Errorf("expected title %s, got %s", trip.Title, got.Title)
	}

	// --- UpdateTrip ---
	trip.Title = "Updated Trip"
	err = srv.UpdateTrip(ctx, trip)
	if err != nil {
		t.Fatalf("UpdateTrip failed: %v", err)
	}
	updated, _ := srv.GetTripByID(ctx, trip.ID)
	if updated.Title != "Updated Trip" {
		t.Errorf("expected updated title, got %s", updated.Title)
	}

	// --- DeleteTrip ---
	err = srv.DeleteTrip(ctx, trip.ID)
	if err != nil {
		t.Fatalf("DeleteTrip failed: %v", err)
	}
	_, err = srv.GetTripByID(ctx, trip.ID)
	if err == nil {
		t.Fatal("expected error for deleted trip")
	}
}
