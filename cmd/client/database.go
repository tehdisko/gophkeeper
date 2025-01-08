package main

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

func getPassword(db *sql.DB, id int) (Password, error) {
	query := `SELECT id, site, username, password FROM passwords WHERE id = ?`
	var p Password
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Site, &p.UserName, &p.Password)
	if err != nil {
		return Password{}, err
	}
	return p, nil
}

func addPassword(db *sql.DB, site, username, password string) (int64, error) {
	query := `INSERT INTO passwords (site, username, password) VALUES (?, ?, ?)`
	result, err := db.Exec(query, site, username, password)
	if err != nil {
		return 0, err
	}

	// Получение ID добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func updatePassword(db *sql.DB, id int, site, username, password string) error {
	query := `UPDATE passwords SET site = ?, username = ?, password = ? WHERE id = ?`
	_, err := db.Exec(query, site, username, password, id)
	return err
}

func deletePassword(db *sql.DB, id int) error {
	query := `DELETE FROM passwords WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func getAllPasswords(db *sql.DB) ([]Password, error) {
	query := `SELECT id, site, username, password FROM passwords`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []Password
	for rows.Next() {
		var p Password
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
