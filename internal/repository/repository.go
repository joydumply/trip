package repository

import (
	"context"
	"fmt"
	"trip-planner/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	GetTrips(ctx context.Context, title string, leaderID *int64, limit, offset int, sortBy, order string) ([]model.Trip, error)
	GetTripByID(ctx context.Context, id int64) (*model.Trip, error)
	UpdateTrip(ctx context.Context, trip *model.Trip) error
	DeleteTrip(ctx context.Context, id int64) error
	Close()
}

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgressRepo(dbURL string) (*PostgresRepo, error) {
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	repo := &PostgresRepo{db: db}

	err = repo.ensureTables()
	if err != nil {
		return nil, err
	}
	return repo, nil

}

func (r *PostgresRepo) ensureTables() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS trips (
		id SERIAL PRIMARY KEY,
		leader_id BIGINT NOT NULL,
		title TEXT NOT NULL,
		description TEXT
	);
	`

	_, err := r.db.Exec(context.Background(), createTableQuery)
	return err
}

func (r *PostgresRepo) CreateTrip(ctx context.Context, trip *model.Trip) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO trips (leader_id, title, description) VALUES ($1, $2, $3) RETURNING id`,
		trip.LeaderID, trip.Title, trip.Description).Scan(&trip.ID)
	if err != nil {
		fmt.Println("CreateTrip error:", err)
	}
	return err
}

func (r *PostgresRepo) GetTrips(ctx context.Context, title string, leaderID *int64, limit, offset int, sortBy, order string) ([]model.Trip, error) {
	query := `SELECT id, leader_id, title, description FROM trips WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argID)
		args = append(args, "%"+title+"%")
		argID++
	}

	if leaderID != nil {
		query += fmt.Sprintf(" AND leader_id = $%d", argID)
		args = append(args, *leaderID)
		argID++

	}

	if sortBy != "id" && sortBy != "title" && sortBy != "leader_id" {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []model.Trip
	for rows.Next() {
		var trip model.Trip
		err = rows.Scan(&trip.ID, &trip.LeaderID, &trip.Title, &trip.Description)
		if err != nil {
			fmt.Println("GetTrips scan error:", err)
			return nil, err
		}
		trips = append(trips, trip)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("GetTrips rows error:", err)
		return nil, err
	}

	return trips, nil
}

func (r *PostgresRepo) GetTripByID(ctx context.Context, id int64) (*model.Trip, error) {
	query := `SELECT id, leader_id, title, description FROM trips WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var trip model.Trip
	err := row.Scan(&trip.ID, &trip.LeaderID, &trip.Title, &trip.Description)
	if err != nil {
		fmt.Println("GetTripByID scan error:", err)

		return nil, err
	}
	return &trip, nil
}

func (r *PostgresRepo) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	query := `UPDATE trips SET leader_id = $1, title = $2, description = $3 WHERE id = $4`
	result, err := r.db.Exec(context.Background(), query, trip.LeaderID, trip.Title, trip.Description, trip.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("trip with id %d not found", trip.ID)
	}
	return nil
}

func (r *PostgresRepo) DeleteTrip(ctx context.Context, id int64) error {
	query := `DELETE FROM trips WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("trip with id %d not found", id)
	}
	return nil
}

func (r *PostgresRepo) Close() {
	if r.db != nil {
		r.db.Close()
	}
}
