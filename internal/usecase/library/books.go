package library

import (
	"context"

	"github.com/project/library/internal/entity"
	"go.uber.org/zap"
)

func (l *impl) RegisterBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	createdBook, err := l.booksRepository.CreateBook(ctx, book)
	if err != nil {
		l.logger.Error("Failed to register book",
			zap.String("bookName", book.Name),
			zap.Error(err))
		return entity.Book{}, err
	}

	return createdBook, nil
}

func (l *impl) GetBook(ctx context.Context, bookID string) (entity.Book, error) {
	book, err := l.booksRepository.GetBook(ctx, bookID)
	if err != nil {
		l.logger.Error("Failed to get book",
			zap.String("bookID", bookID),
			zap.Error(err))
		return entity.Book{}, err
	}

	return book, nil
}

func (l *impl) UpdateBook(ctx context.Context, bookID string, name string, authorIDs []string) error {
	updatedBook := entity.Book{
		ID:        bookID,
		Name:      name,
		AuthorIDs: authorIDs,
	}

	if err := l.booksRepository.UpdateBook(ctx, updatedBook); err != nil {
		l.logger.Error("Failed to update book in repository",
			zap.String("bookID", bookID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (l *impl) GetAuthorBooks(ctx context.Context, authorID string) ([]entity.Book, error) {
	books, err := l.booksRepository.GetBooksByAuthor(ctx, authorID)
	if err != nil {
		l.logger.Error("Failed to get books by author",
			zap.String("authorID", authorID),
			zap.Error(err))
		return nil, err
	}

	return books, nil
}
