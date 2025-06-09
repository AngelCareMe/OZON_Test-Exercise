package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"ozon_test/config"
	"ozon_test/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrCommentsNotAllowed = errors.New("comments not allowed")

type PostStorage interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	UpdatePost(post *models.Post) error
}

type CommentStorage interface {
	CreateComment(comment *models.Comment) error
	GetCommentsByPostID(postID int, limit, offset int) ([]*models.Comment, error)
}

type InMemoryPostStorage struct {
	posts  map[int]*models.Post
	mu     sync.Mutex
	nextID int
}

type InMemoryCommentStorage struct {
	comments map[int]*models.Comment
	mu       sync.Mutex
	nextID   int
}

func NewInMemoryPostStorage() *InMemoryPostStorage {
	return &InMemoryPostStorage{
		posts:  make(map[int]*models.Post),
		nextID: 1,
	}
}

func NewInMemoryCommentStorage() *InMemoryCommentStorage {
	return &InMemoryCommentStorage{
		comments: make(map[int]*models.Comment),
		nextID:   1,
	}
}

func (s *InMemoryPostStorage) CreatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	post.ID = s.nextID
	s.nextID++
	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryPostStorage) GetPostByID(id int) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	post, exists := s.posts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return post, nil
}

func (s *InMemoryPostStorage) GetAllPosts() ([]*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *InMemoryPostStorage) UpdatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.posts[post.ID]; !exists {
		return ErrNotFound
	}
	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryCommentStorage) CreateComment(comment *models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	comment.ID = s.nextID
	s.nextID++
	s.comments[comment.ID] = comment
	return nil
}

func (s *InMemoryCommentStorage) GetCommentsByPostID(postID int, limit, offset int) ([]*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var comments []*models.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	start := offset
	if start >= len(comments) {
		return []*models.Comment{}, nil
	}
	end := start + limit
	if end > len(comments) {
		end = len(comments)
	}
	return comments[start:end], nil
}

type PostgresPostStorage struct {
	pool *pgxpool.Pool
}

type PostgresCommentStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresPostStorage(pool *pgxpool.Pool) *PostgresPostStorage {
	return &PostgresPostStorage{pool: pool}
}

func NewPostgresCommentStorage(pool *pgxpool.Pool) *PostgresCommentStorage {
	return &PostgresCommentStorage{pool: pool}
}

func (s *PostgresPostStorage) CreatePost(post *models.Post) error {
	query := squirrel.Insert("posts").Columns("title", "text", "allow_comments", "author", "created_at").
		Values(post.Title, post.Text, post.AllowComments, post.Author, post.CreatedAt).
		Suffix("RETURNING id").PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = s.pool.QueryRow(context.Background(), sql, args...).Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresPostStorage) GetPostByID(id int) (*models.Post, error) {
	query := squirrel.Select("id", "title", "text", "allow_comments", "author", "created_at").
		From("posts").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := s.pool.QueryRow(context.Background(), sql, args...)
	post := &models.Post{}
	err = row.Scan(&post.ID, &post.Title, &post.Text, &post.AllowComments, &post.Author, &post.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}

func (s *PostgresPostStorage) GetAllPosts() ([]*models.Post, error) {
	query := squirrel.Select("id", "title", "text", "allow_comments", "author", "created_at").
		From("posts").PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Text, &post.AllowComments, &post.Author, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostgresPostStorage) UpdatePost(post *models.Post) error {
	query := squirrel.Update("posts").Set("title", post.Title).Set("text", post.Text).
		Set("allow_comments", post.AllowComments).Set("author", post.Author).
		Set("created_at", post.CreatedAt).Where(squirrel.Eq{"id": post.ID}).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := s.pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresCommentStorage) CreateComment(comment *models.Comment) error {
	query := squirrel.Insert("comments").Columns("post_id", "parent_comment_id", "text", "author", "created_at").
		Values(comment.PostID, comment.ParentCommentID, comment.Text, comment.Author, comment.CreatedAt).
		Suffix("RETURNING id").PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = s.pool.QueryRow(context.Background(), sql, args...).Scan(&comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresCommentStorage) GetCommentsByPostID(postID int, limit, offset int) ([]*models.Comment, error) {
	query := squirrel.Select("id", "post_id", "parent_comment_id", "text", "author", "created_at").
		From("comments").Where(squirrel.Eq{"post_id": postID}).
		Limit(uint64(limit)).Offset(uint64(offset)).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var parentID *int
		err = rows.Scan(&comment.ID, &comment.PostID, &parentID, &comment.Text, &comment.Author, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.ParentCommentID = parentID
		comments = append(comments, comment)
	}
	return comments, nil
}

func CreateDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора строки подключения: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return pool, nil
}
