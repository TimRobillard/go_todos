package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Todo struct {
	Id       int
	Title    string
	Complete bool
}

type Storage interface {
	CreateTodo(string) (*Todo, error)
	DeleteTodo(int) error
	GetAllTodos() ([]*Todo, error)
	ToggleComplete(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if db, err := sql.Open("postgres", "user=tim dbname=todos password=password sslmode=disable"); err != nil {
		return nil, err
	} else {
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return &PostgresStore{
			db: db,
		}, nil
	}
}

func (s *PostgresStore) Init() error {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		complete BOOL DEFAULT false
	);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateToDo(title string) (*Todo, error) {
	query := `INSERT INTO todo (title) 
	VALUES($1)
	RETURNING id;`
	fmt.Printf("title: %s", title)

	var id int
	err := s.db.QueryRow(query, title).Scan(&id)

	return &Todo{
		Id:       id,
		Title:    title,
		Complete: false,
	}, err

}

func (s *PostgresStore) DeleteToDo(id int) error {
	query := `DELETE todo WHERE id = $1;`

	_, err := s.db.Exec(query, id)

	return err
}

func (s *PostgresStore) GetAllTodos() ([]*Todo, error) {
	query := `SELECT * FROM todo ORDER BY id;`

	var todos []*Todo
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("1")
	for rows.Next() {
		fmt.Println("2-1")
		todo := &Todo{} 
		fmt.Println("2")
		rows.Scan(&todo.Id, &todo.Title, &todo.Complete)
		fmt.Println("3")
		todos = append(todos, todo)
	}

	return todos, err

}

func (s *PostgresStore) ToggleTodo(id int) error {
	query := `UPDATE todo 
	SET complete = NOT complete
	WHERE id = $1;`

	_, err := s.db.Exec(query, id)

	return err
}
