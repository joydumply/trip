package repository

import (
	"context"
	"testing"
	"trip-planner/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestRepo(t *testing.T) *PostgresRepo {
	dbURL := "postgres://postgres:postgres@localhost:5432/trip_planner?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	// создаём временную таблицу для тестов
	createTable := `
	CREATE TABLE IF NOT EXISTS trips_test (
		id SERIAL PRIMARY KEY,
		leader_id BIGINT NOT NULL,
		title TEXT NOT NULL,
		description TEXT
	);
	`
	_, err = pool.Exec(context.Background(), createTable)
	if err != nil {
		t.Fatalf("failed to create trips_test table: %v", err)
	}

	return &PostgresRepo{db: pool}
}

func teardownTestRepo(t *testing.T, r *PostgresRepo) {
	_, err := r.db.Exec(context.Background(), "DROP TABLE IF EXISTS trips_test")
	if err != nil {
		t.Fatalf("failed to drop trips_test table: %v", err)
	}
	r.Close()
}

func TestRepositoryCRUD(t *testing.T) {
	r := setupTestRepo(t)
	defer teardownTestRepo(t, r)

	ctx := context.Background()

	// CREATE
	trip := &model.Trip{
		LeaderID:    1,
		Title:       "Test trip",
		Description: "Test description",
	}
	createQuery := `INSERT INTO trips_test (leader_id, title, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, createQuery, trip.LeaderID, trip.Title, trip.Description).Scan(&trip.ID)
	if err != nil {
		t.Fatalf("failed to create trip: %v", err)
	}

	// READ
	var got model.Trip
	readQuery := `SELECT id, leader_id, title, description FROM trips_test WHERE id = $1`
	err = r.db.QueryRow(ctx, readQuery, trip.ID).Scan(&got.ID, &got.LeaderID, &got.Title, &got.Description)
	if err != nil {
		t.Fatalf("failed to get trip: %v", err)
	}
	if got.Title != trip.Title {
		t.Fatalf("expected title %s, got %s", trip.Title, got.Title)
	}

	// UPDATE
	trip.Title = "Updated trip"
	updateQuery := `UPDATE trips_test SET title = $1 WHERE id = $2`
	_, err = r.db.Exec(ctx, updateQuery, trip.Title, trip.ID)
	if err != nil {
		t.Fatalf("failed to update trip: %v", err)
	}

	// DELETE
	deleteQuery := `DELETE FROM trips_test WHERE id = $1`
	_, err = r.db.Exec(ctx, deleteQuery, trip.ID)
	if err != nil {
		t.Fatalf("failed to delete trip: %v", err)
	}
}
