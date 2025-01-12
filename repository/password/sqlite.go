package password

import (
	"database/sql"
	"fmt"
	"log"
)
import _ "github.com/mattn/go-sqlite3"
import "github.com/tehdisko/gophkeeper/domain"

type Repository struct {
	db *sql.DB
}

func New(dbName string) *Repository {
	db, err := sql.Open("sqlite3", fmt.Sprintf("db/%s.db", dbName))
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPassword(id int) (domain.Password, error) {
	query := `SELECT id, site, username, password FROM passwords WHERE id = ?`
	var p domain.Password
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Site, &p.UserName, &p.Password)
	if err != nil {
		return domain.Password{}, err
	}
	return p, nil
}

func (r *Repository) AddPassword(site, username, password string) (int64, error) {
	query := `INSERT INTO passwords (site, username, password) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, site, username, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdatePassword(id int, site, username, password string) error {
	query := `UPDATE passwords SET site = ?, username = ?, password = ? WHERE id = ?`
	_, err := r.db.Exec(query, site, username, password, id)

	return err
}

func (r *Repository) DeletePassword(id int) error {
	query := `DELETE FROM passwords WHERE id = ?`
	_, err := r.db.Exec(query, id)

	return err
}

func (r *Repository) GetAllPasswords() ([]domain.Password, error) {
	query := `SELECT id, site, username, password FROM passwords`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var passwords []domain.Password

	for rows.Next() {
		var p domain.Password
		err := rows.Scan(&p.ID, &p.Site, &p.UserName, &p.Password)
		if err != nil {
			return nil, err
		}
		passwords = append(passwords, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}
