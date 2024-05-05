package storage

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Для хранения пароля
var pwd string = os.Getenv("psqpwd")

// Для хранения пути подключения и передачи в функцию NewConn
var dbString string = "postgres://postgres:" + pwd + "@localhost/test_tasks"

type Storage struct {
	db *pgxpool.Pool
}

func NewConn(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(),
		constr)
	if err != nil {
		log.Fatalf("Unable to connect to database. Error: /%v/n", err)
	}
	s := Storage{db: db}
	return &s, nil
}

type Task struct {
	Id          int
	Opened      int
	Closed      int
	Author_id   int
	Assigned_id int
	Title       string
	Content     string
}

// Функция getAllTasks возваращает все текущие задачи
func (s *Storage) getAllTasks() ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT * FROM tasks;
	`)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.Id,
			&t.Opened,
			&t.Closed,
			&t.Author_id,
			&t.Assigned_id,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()

}

// Функция newTasks принимает параметры новой задачи
// и заносит их в базу данных
func (s *Storage) newTask(author int,
	assigned int, title string, content string) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO tasks(author_id, assigned id, title, content) VALUES ($1,$2,$3,$4);
	`, author, assigned, title, content)
	if err != nil {
		return err
	}
	return nil
}

// Функция tasksByAuthor принимает id автора и возвращает все его задачи
func (s *Storage) tasksByAuthor(author_id int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT * FROM tasks WHERE author_id = ($1);
	`, author_id)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.Id,
			&t.Opened,
			&t.Closed,
			&t.Author_id,
			&t.Assigned_id,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()

}

// Функция tasksByLabel принимает имя ярлыка и вовзвращает список задач,
// относящихся к этому ярлыку
func (s *Storage) tasksByLabel(label string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT tasks.id, tasks.opened, tasks.closed, tasks.author_id, tasks.assigned_id, tasks.title, tasks.content 
	FROM tasks, tasks_lables, labels
	WHERE tasks.id=tasks_lables.task_id
	AND tasks_lables.label_id=(SELECT labels.id FROM labels WHERE labels.name=$1);`, label)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.Id,
			&t.Opened,
			&t.Closed,
			&t.Author_id,
			&t.Assigned_id,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// Функция updateTask принимает id и дату закрытия задачи
// и вносит изменения в БД (если я правильно понял задачу)
func (s *Storage) updateTask(task_id, closeDate int) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE tasks
	SET closed =$2
	WHERE id=$1;
	`, task_id, closeDate)
	if err != nil {
		return err
	}
	return nil
}

// Функция deleteTask принимает id задачи удаляет ее
func (s *Storage) deleteTask(task_id int) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks
	WHERE id=$1;
	`, task_id)
	if err != nil {
		return err
	}
	return nil
}
