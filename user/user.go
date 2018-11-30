package user

import (
	"database/sql"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BankAccount struct {
	ID      int    `json:"id"`
	userID  int    `json:"user_id"`
	Number  int    `json:"number"`
	Name    string `json:"name"`
	balance int    `json:"balance"`
}

type Service struct {
	DB *sql.DB
}

func (s *Service) FindByID(id int) (*User, error) {
	stmt := "SELECT id, first_name, last_name, email FROM users WHERE id = $1"
	row := s.DB.QueryRow(stmt, id)
	var u User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Insert(u *User) error {
	stmt := `INSERT INTO users(first_name, last_name, email)
		 values ($1, $2, $3) RETURNING id`
	row := s.DB.QueryRow(stmt, u.FirstName, u.LastName)
	err := row.Scan(&u.ID)

	return err
}

func (s *Service) All() ([]User, error) {
	stmt := "SELECT id, first_name, last_name, email FROM users ORDER BY id DESC"
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	var us []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

func (s *Service) Update(u *User) error {
	stmt := "UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4"
	_, err := s.DB.Exec(stmt, u.FirstName, u.LastName, u.ID)
	return err
}

func (s *Service) Delete(u *User) error {
	stmt := "DELETE FROM users WHERE id = $1"
	_, err := s.DB.Exec(stmt, u.ID)
	return err
}
