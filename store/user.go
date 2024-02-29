package store

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type SafeUser struct {
	Id       int
	Username string
}

type User struct {
	Id                int
	Username          string
	EncryptedPassword string
}

type UserStorage interface {
	InitUser() error
	CreateUser(string, string) (*SafeUser, error)
	GetUserById(int) (*User, error)
	GetUserByUsername(string) (*User, error)
	DeleteUser(int) error
}

func (pg PostgresStore) InitUser() error {
	fmt.Println("creating users table")
	query := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255),
		is_deleted BOOL DEFAULT false
	);`

	_, err := pg.db.Exec(query)
	return err
}

func (pg PostgresStore) CreateUser(u, password string) (*SafeUser, error) {
	username := strings.ToLower(u)
	username = strings.Trim(username, " ")

	query := `INSERT INTO users(username, password)
	VALUES($1, $2)
	RETURNING id;`

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	var id int
	err = pg.db.QueryRow(query, username, string(bytes)).Scan(&id)

	user := &SafeUser{
		Id:       id,
		Username: username,
	}
	return user, err
}

func (pg PostgresStore) GetUserById(id int) (*User, error) {
	query := `SELECT id, username, password
	FROM users WHERE id = $1 AND is_deleted = false;`

	user := new(User)
	err := pg.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.EncryptedPassword)

	return user, err
}

func (pg PostgresStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password
	FROM users WHERE username = $1 AND is_deleted = false;`

	user := new(User)
	err := pg.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.EncryptedPassword)

	return user, err
}

func (pg PostgresStore) DeleteUser(id int) error {
	query := `UPDATE users SET is_deleted = true WHERE id = $1;`

	_, err := pg.db.Exec(query, id)

	return err
}

func (u User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}
