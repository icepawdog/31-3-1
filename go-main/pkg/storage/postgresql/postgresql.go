package postgresql

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Tasks(taskID, authorID int) ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			author_id,
			title,
			content,
			created_at
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, p)

	}
	return tasks, rows.Err()
}

func (s *Storage) NewTask(p storage.Post) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO post (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		p.Title,
		p.Content,
	).Scan(&id)
	return id, err
}

func (s *Storage) DelTask(p storage.Post) error {
	_, err := s.db.Exec(context.Background(),
		`DELETE FROM tasks WHERE id = $1;`,
		p.ID,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Storage) UpdateTask(p storage.Post) error {
	_, err := s.db.Exec(context.Background(),
		`UPDATE tasks SET title = $1, content = $2 WHERE id = $3;`,
		p.ID,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
