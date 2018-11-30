package user

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BankAccount struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Number  int    `json:"number"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type Service struct {
	DB *sql.DB
}

func (s *Service) FindByID(id int) (*User, error) {
	stmt := "SELECT id, first_name, last_name FROM users WHERE id = $1"
	row := s.DB.QueryRow(stmt, id)
	var u User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Insert(u *User) error {
	stmt := `INSERT INTO Users(first_name, last_name)
		 values ($1, $2) RETURNING id`
	row := s.DB.QueryRow(stmt, u.FirstName, u.LastName)
	err := row.Scan(&u.ID)

	return err
}

func (s *Service) All() ([]User, error) {
	stmt := "SELECT id, first_name, last_name FROM Users ORDER BY id DESC"
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
	stmt := "UPDATE Users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4"
	_, err := s.DB.Exec(stmt, u.FirstName, u.LastName, u.ID)
	return err
}

func (s *Service) Delete(u *User) error {
	stmt := "DELETE FROM Users WHERE id = $1"
	_, err := s.DB.Exec(stmt, u.ID)
	return err
}

//BankAccount Service
func (s *Service) AddBankAc(bkAc *BankAccount) error {
	stmt := `INSERT INTO BankAccount(user_id, number,name,balance)
		 values ($1, $2,$3,$4) RETURNING id`
	fmt.Println(bkAc.UserID)
	row := s.DB.QueryRow(stmt, bkAc.UserID, bkAc.Number, bkAc.Name, bkAc.Balance)
	err := row.Scan(&bkAc.ID)

	return err
}

func (s *Service) GetAllUserBkAc(id int) ([]BankAccount, error) {
	stmt, _ := s.DB.Prepare("SELECT * FROM BankAccount WHERE user_id = $1")
	rows, _ := stmt.Query(id)

	var bkAc []BankAccount
	for rows.Next() {
		var bk BankAccount
		err := rows.Scan(&bk.ID, &bk.UserID, &bk.Number, &bk.Name, &bk.Balance)
		if err != nil {
			return nil, err
		}
		bkAc = append(bkAc, bk)
	}

	return bkAc, nil
}

// func (s *Service) RemoveBkAc(u *User) error {
// 	stmt := "DELETE FROM users WHERE id = $1"
// 	_, err := s.DB.Exec(stmt, u.ID)
// 	return err
// }
// func (s *Service) Withdraw(u *User) error {
// 	stmt := "DELETE FROM users WHERE id = $1"
// 	_, err := s.DB.Exec(stmt, u.ID)
// 	return err
// }
// func (s *Service) Deposit(u *User) error {
// 	stmt := "DELETE FROM users WHERE id = $1"
// 	_, err := s.DB.Exec(stmt, u.ID)
// 	return err
// }
// func (s *Service) Transfers(u *User) error {
// 	stmt := "DELETE FROM users WHERE id = $1"
// 	_, err := s.DB.Exec(stmt, u.ID)
// 	return err
// }
