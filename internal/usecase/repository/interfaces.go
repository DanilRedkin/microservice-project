package repository

//go:generate mockgen -source=interfaces.go -destination=repository_mock.go -package=repository

import (
	"context"

	"github.com/project/library/internal/entity"
)

type (
	AuthorRepository interface {
		CreateAuthor(ctx context.Context, author entity.Author) (entity.Author, error)
		GetAuthor(ctx context.Context, authorID string) (entity.Author, error)
		ChangeAuthor(ctx context.Context, author entity.Author) error
	}

	BooksRepository interface {
		CreateBook(ctx context.Context, book entity.Book) (entity.Book, error)
		UpdateBook(ctx context.Context, book entity.Book) error
		GetBook(ctx context.Context, bookID string) (entity.Book, error)
		GetBooksByAuthor(ctx context.Context, authorID string) ([]entity.Book, error)
	}
)
