package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}
func (s *PostgresStore) init() error {
	err := s.CreateAccountTable()
	return err
}
func (s *PostgresStore) CreateAccountTable() error {
	//!https://golangbot.com/mysql-create-table-insert-row/
	query := `CREATE TABLE IF NOT EXISTS account(
		id serial primary key , 
		firstName text, 
		lastName text, 
		number serial, 
		balance serial,
		created_at TIMESTAMP default CURRENT_TIMESTAMP
		)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelfunc()
	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating account table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func (s *PostgresStore) CreateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(int) (*Account, error) {
	return &Account{}, nil
}
