package repository

import (
	"database/sql"
	"go-fiber-crud/internal/models"
)

type TodoRepository struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{DB: db}
}

func (r *TodoRepository) GetTodos(userId string) ([]models.Todo, error) {
	rows, err := r.DB.Query("SELECT id, title FROM todos WHERE user_Id = $1 ORDER BY id DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			continue
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (r *TodoRepository) CreateTodo(userId, title string) error {
	_, err := r.DB.Exec("INSERT INTO todos (user_id, title) VALUES ($1, $2)", userId, title)
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepository) DeleteTodo(todoId string) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id = $1", todoId)
	if err != nil {
		return err
	}

	return nil
}
