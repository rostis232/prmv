package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rostis232/prmv/models"
)

const (
	postsTable = "posts"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(configDB string) (*Postgres, error) {
	db, err := sqlx.Open("postgres", configDB)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	p := Postgres{db: db}

	return &p, nil
}

func (p *Postgres) AddPost(post models.Post) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (title, content) values ($1, $2) returning id", postsTable)
	err := p.db.QueryRow(query, post.Title, post.Content).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error adding post: %w", err)
	}
	return id, nil
}

func (p *Postgres) GetAllPosts() ([]models.Post, error) {
	var posts = []models.Post{}
	query := fmt.Sprintf("select * from %s", postsTable)
	err := p.db.Select(&posts, query)
	if err != nil {
		return posts, fmt.Errorf("error getting all posts: %w", err)
	}
	return posts, nil
}
func (p *Postgres) UpdatePost(post models.Post) (int, error) {
	var id int
	query := fmt.Sprintf("update %s set title = $1, content = $2 where id = $3 returning id", postsTable)
	err := p.db.QueryRow(query, post.Title, post.Content, post.ID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error updating post: %w", err)
	}
	return id, nil
}
func (p *Postgres) GetPost(id int) (models.Post, error) {
	var post models.Post
	query := fmt.Sprintf("select * from %s where id = $1", postsTable)
	err := p.db.Get(&post, query, id)
	if err != nil {
		return post, fmt.Errorf("error getting post: %w", err)
	}
	return post, nil
}
func (p *Postgres) DeletePost(id int) error {
	query := fmt.Sprintf("delete from %s where id = $1", postsTable)
	_, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}
