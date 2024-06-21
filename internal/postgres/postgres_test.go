package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rostis232/prmv/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testDB = "port=5434 user=gopher password=some_pass dbname=postsdb sslmode=disable"
)

func TestNewPostgres(t *testing.T) {
	_, err := NewPostgres(testDB)
	if err != nil {
		t.Error(err)
	}
}

func prepareTestDB() (*Postgres, error) {
	p, err := NewPostgres(testDB)
	if err != nil {
		return nil, err
	}

	createQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`, postsTable)
	_, err = p.db.Exec(createQuery)
	if err != nil {
		return nil, err
	}

	truncateQuery := fmt.Sprintf(`TRUNCATE TABLE %s`, postsTable)
	_, err = p.db.Exec(truncateQuery)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func TestAddPost(t *testing.T) {
	p, err := prepareTestDB()
	if err != nil {
		t.Error(err)
	}

	post := models.Post{Title: "Test Title", Content: "Test Content"}
	id, err := p.AddPost(post)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	var count int
	query := fmt.Sprintf("select count(*) from %s where id=$1", postsTable)
	err = p.db.Get(&count, query, id)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetAllPosts(t *testing.T) {
	p, err := prepareTestDB()
	if err != nil {
		t.Error(err)
	}

	posts := []models.Post{
		{Title: "Test Title 1", Content: "Test Content 1"},
		{Title: "Test Title 2", Content: "Test Content 2"},
	}

	for _, post := range posts {
		_, err := p.AddPost(post)
		assert.NoError(t, err)
	}

	allPosts, err := p.GetAllPosts()
	assert.NoError(t, err)
	assert.Len(t, allPosts, len(posts))
}

func TestUpdatePost(t *testing.T) {
	p, err := prepareTestDB()
	if err != nil {
		t.Error(err)
	}

	post := models.Post{Title: "Original Title", Content: "Original Content"}
	id, err := p.AddPost(post)
	assert.NoError(t, err)

	post.ID = id
	post.Title = "Updated Title"
	post.Content = "Updated Content"
	updatedID, err := p.UpdatePost(post)
	assert.NoError(t, err)
	assert.Equal(t, id, updatedID)

	updatedPost, err := p.GetPost(id)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, updatedPost.Title)
	assert.Equal(t, post.Content, updatedPost.Content)
}

func TestGetPost(t *testing.T) {
	p, err := prepareTestDB()
	if err != nil {
		t.Error(err)
	}

	post := models.Post{Title: "Test Title", Content: "Test Content"}
	id, err := p.AddPost(post)
	assert.NoError(t, err)

	fetchedPost, err := p.GetPost(id)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, fetchedPost.Title)
	assert.Equal(t, post.Content, fetchedPost.Content)
}

func TestDeletePost(t *testing.T) {
	p, err := prepareTestDB()
	if err != nil {
		t.Error(err)
	}

	post := models.Post{Title: "Test Title", Content: "Test Content"}
	id, err := p.AddPost(post)
	assert.NoError(t, err)

	err = p.DeletePost(id)
	assert.NoError(t, err)

	var count int
	query := fmt.Sprintf("select count(*) from %s where id=$1", postsTable)
	err = p.db.Get(&count, query, id)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
