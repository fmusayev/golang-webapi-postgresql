package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	models "../model"
)

// NewPostRepository ...
func NewPostRepository(Conn *sql.DB) PostRepository {
	return PostRepository{
		Conn: Conn,
	}
}

// PostRepository ...
type PostRepository struct {
	Conn *sql.DB
}

func (postRepo *PostRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := postRepo.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Fatalln("error getting rows", err)
		return nil, err
	}
	defer rows.Close()

	payload := make([]*models.Post, 0)
	for rows.Next() {
		data := new(models.Post)

		err := rows.Scan(
			&data.ID,
			&data.Title,
		)
		if err != nil {
			log.Fatalln("error scanning row", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

// Fetch ...
func (postRepo *PostRepository) Fetch(ctx context.Context, num int64) ([]*models.Post, error) {
	query := "select id, title from posts"
	return postRepo.fetch(ctx, query)
}

// GetByID ...
func (postRepo *PostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := "select id, title from posts where id = $1"

	rows, err := postRepo.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &models.Post{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("repo.post.getById.not-found")
	}

	return payload, nil
}

// Create ..,
func (postRepo *PostRepository) Create(ctx context.Context, p *models.Post) (int64, error) {
	query := "insert into posts(title) values($1)"

	stmt, err := postRepo.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln("error preparing ctx:", err)
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Title)
	defer stmt.Close()

	if err != nil {
		log.Fatalln("error executing context: ", err)
		return -1, err
	}

	return res.RowsAffected()
}

// Update ...
func (postRepo *PostRepository) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
	query := "Update posts set title=$1 where id=$2"

	stmt, err := postRepo.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		p.Title,
		p.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return p, nil
}

// Delete ...
func (postRepo *PostRepository) Delete(ctx context.Context, id int64) (bool, error) {
	query := "Delete From posts Where id=?"

	stmt, err := postRepo.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
