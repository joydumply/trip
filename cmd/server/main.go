package main

import (
	"context"
	"database/sql"
	"log"

	"trip-planner/internal/app"
	"trip-planner/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx как database/sql драйвер
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// --- Подключение для приложения через pgxpool ---
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// --- Подключение для golang-migrate через database/sql ---
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// --- Применяем миграции ---
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

    // ! If you want to roll back migrations, uncomment the following lines

	// if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal(err)
	// }
	// log.Println("Migrations rolled back (trips table dropped)")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	log.Println("Migrations applied successfully")

	// --- Запуск приложения ---
	app.Run(cfg)
}
