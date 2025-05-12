package postgres

import "database/sql"

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(connectionStr string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	return &PostgresRepo{db}, nil
}
