package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"ozon_test/config"
)

func ApplyMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Ожидание готовности базы данных
	if err := waitForDB(dsn); err != nil {
		return fmt.Errorf("ошибка ожидания базы данных: %w", err)
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("ошибка создания миграции: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	log.Println("Миграции успешно применены")
	return nil
}

func waitForDB(dsn string) error {
	const maxAttempts = 10
	const delay = 2 * time.Second

	for i := 0; i < maxAttempts; i++ {
		pool, err := pgxpool.Connect(context.Background(), dsn)
		if err == nil {
			pool.Close()
			return nil
		}
		log.Printf("База данных недоступна, попытка %d/%d: %v", i+1, maxAttempts, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("база данных недоступна после %d попыток", maxAttempts)
}
