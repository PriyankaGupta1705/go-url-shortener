package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) RunMigrations() error {
	query, err := os.ReadFile("db/migrations/001_create_urls.sql")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(query))
	return err
}

func (s *PostgresStore) Save(code, originalURL string) error {
	_, err := s.db.Exec(
		"INSERT INTO urls (code, original) VALUES ($1, $2)",
		code, originalURL,
	)
	return err
}

func (s *PostgresStore) Get(code string) (string, error) {
	var original string
	err := s.db.QueryRow(
		"SELECT original FROM urls WHERE code = $1", code,
	).Scan(&original)

	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("url not found")
	}
	return original, err
}

func (s *PostgresStore) IncrementVisits(code string) {
	s.db.Exec("UPDATE urls SET visits = visits + 1 WHERE code = $1", code)
}

func (s *PostgresStore) GetStats(code string) (string, int, error) {
	var original string
	var visits int
	err := s.db.QueryRow(
		"SELECT original, visits FROM urls WHERE code = $1", code,
	).Scan(&original, &visits)

	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, errors.New("url not found")
	}
	return original, visits, err
}
