package repository

import (
	"database/sql"
	"todo-api/model"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) model.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]model.Todo, error) {
	rows, err := r.db.Query("SELECT id, text, done FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (r *todoRepository) Create(text string) (*model.Todo, error) {
	res, err := r.db.Exec("INSERT INTO todos (text) VALUES (?)", text)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &model.Todo{ID: int(id), Text: text, Done: false}, nil
}

func (r *todoRepository) Update(id int, done bool) (*model.Todo, error) {
	if _, err := r.db.Exec("UPDATE todos SET done = ? WHERE id = ?", done, id); err != nil {
		return nil, err
	}
	var t model.Todo
	err := r.db.QueryRow("SELECT id, text, done FROM todos WHERE id = ?", id).Scan(&t.ID, &t.Text, &t.Done)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *todoRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
