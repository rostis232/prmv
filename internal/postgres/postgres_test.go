package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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

	testCases := []struct {
		post          models.Post
		errorExpected bool
	}{
		{
			post:          models.Post{Title: "Test Title", Content: "Test Content"},
			errorExpected: false,
		},
		{
			errorExpected: false,
		},
	}

	for i, tc := range testCases {
		id, err := p.AddPost(tc.post)
		if tc.errorExpected {
			assert.Error(t, err, fmt.Sprintf("case %d", i))
		} else {
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
			assert.NotZero(t, id, fmt.Sprintf("case %d", i))

			var count int
			query := fmt.Sprintf("select count(*) from %s where id=$1", postsTable)
			err = p.db.Get(&count, query, id)
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
			assert.Equal(t, 1, count, fmt.Sprintf("case %d", i))
		}
	}
}

func TestGetAllPosts(t *testing.T) {
	testCases := []struct {
		posts         []models.Post
		errorExpected bool
	}{
		{
			posts: []models.Post{
				{Title: "Test Title 1", Content: "Test Content 1"},
				{Title: "Test Title 2", Content: "Test Content 2"},
			},
			errorExpected: false,
		},
		{
			posts:         []models.Post{},
			errorExpected: false,
		},
	}

	for i, tc := range testCases {
		p, err := prepareTestDB()
		if err != nil {
			t.Error(err)
		}

		for _, post := range tc.posts {
			_, err := p.AddPost(post)
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
		}

		allPosts, err := p.GetAllPosts()

		if tc.errorExpected {
			assert.Error(t, err, fmt.Sprintf("case %d", i))
		} else {
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
		}
		assert.Len(t, allPosts, len(tc.posts), fmt.Sprintf("case %d", i))
	}
}

func TestUpdatePost(t *testing.T) {
	testCases := []struct {
		post           models.Post
		updatedTitle   string
		updatedContent string
		errorExpected  error
	}{
		{
			post:           models.Post{Title: "Original Title", Content: "Original Content"},
			updatedTitle:   "Updated Title",
			updatedContent: "Updated Content",
			errorExpected:  nil,
		},
		{
			errorExpected: sql.ErrNoRows,
		},
	}

	for i, tc := range testCases {
		p, err := prepareTestDB()
		if err != nil {
			t.Error(err)
		}

		var id int

		if tc.errorExpected == nil {
			id, err = p.AddPost(tc.post)
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
		}

		tc.post.ID = id
		tc.post.Title = tc.updatedTitle
		tc.post.Content = tc.updatedContent
		updatedID, err := p.UpdatePost(tc.post)

		if tc.errorExpected != nil {
			assert.Error(t, err, fmt.Sprintf("case %d", i))

			errorType := errors.Is(err, tc.errorExpected)
			assert.Equal(t, true, errorType, fmt.Sprintf("case %d", i))
		} else {
			assert.NoError(t, err, fmt.Sprintf("case %d", i))

			assert.Equal(t, id, updatedID, fmt.Sprintf("case %d", i))

			updatedPost, err := p.GetPost(id)
			assert.NoError(t, err, fmt.Sprintf("case %d", i))

			assert.Equal(t, tc.post.Title, updatedPost.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, tc.post.Content, updatedPost.Content, fmt.Sprintf("case %d", i))
		}
	}
}

func TestGetPost(t *testing.T) {
	testCases := []struct {
		post          models.Post
		errorExpected error
	}{
		{
			post:          models.Post{Title: "Test Title", Content: "Test Content"},
			errorExpected: nil,
		},
		{
			errorExpected: sql.ErrNoRows,
		},
	}

	for i, tc := range testCases {
		p, err := prepareTestDB()
		if err != nil {
			t.Error(err)
		}

		var id int

		if tc.errorExpected == nil {
			id, err = p.AddPost(tc.post)
			assert.NoError(t, err, fmt.Sprintf("case %d", i))
		}

		fetchedPost, err := p.GetPost(id)
		if tc.errorExpected == nil {
			assert.NoError(t, err, fmt.Sprintf("case %d", i))

			assert.Equal(t, tc.post.Title, fetchedPost.Title, fmt.Sprintf("case %d", i))
			assert.Equal(t, tc.post.Content, fetchedPost.Content, fmt.Sprintf("case %d", i))
		} else {
			assert.Error(t, err, fmt.Sprintf("case %d", i))

			errorType := errors.Is(err, tc.errorExpected)
			assert.Equal(t, true, errorType, fmt.Sprintf("case %d", i))
		}
	}
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
