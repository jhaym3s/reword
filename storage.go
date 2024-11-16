package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
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

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account (firstName, lastName,number,balance, created_at) VALUES ($1, $2,$3,$4,$5)`
	resp, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", resp)
	return nil
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next(){
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d does not exist",id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, err
}
