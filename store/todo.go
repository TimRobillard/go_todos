package store

import "fmt"

type Todo struct {
	Id       int
	Title    string
	UserId   int
	Complete bool
}

type TodoStorage interface {
	Init() error
	CreateTodo(string, int) (*Todo, error)
	DeleteTodo(int) error
	GetAllTodos(int) ([]*Todo, error)
	ToggleComplete(int, int) error
}

func (s *PostgresStore) Init() error {
	fmt.Print("Creating table todo")
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		complete BOOL DEFAULT false
	);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateTodo(title string, userId int) (*Todo, error) {
	query := `INSERT INTO todo (title, user_id) 
	VALUES($1, $2)
	RETURNING id;`

	var id int
	err := s.db.QueryRow(query, title, userId).Scan(&id)

	return &Todo{
		Id:       id,
		Title:    title,
		UserId:   userId,
		Complete: false,
	}, err

}

func (s *PostgresStore) DeleteTodo(id int) error {
	query := `DELETE todo WHERE id = $1;`

	_, err := s.db.Exec(query, id)

	return err
}

func (s *PostgresStore) GetAllTodos(userId int) ([]*Todo, error) {
	query := `SELECT id, title, complete FROM todo WHERE user_id = $1 ORDER BY id;`

	var todos []*Todo
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		todo := &Todo{}
		rows.Scan(&todo.Id, &todo.Title, &todo.Complete)
		todos = append(todos, todo)
	}

	return todos, err

}

func (s *PostgresStore) ToggleComplete(id int, userId int) error {
	query := `UPDATE todo 
	SET complete = NOT complete
	WHERE id = $1 AND user_id = $2;`

	_, err := s.db.Exec(query, id, userId)

	return err
}
