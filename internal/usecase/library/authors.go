package library

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	"go.uber.org/zap"
)

func (l *impl) RegisterAuthor(ctx context.Context, authorName string) (entity.Author, error) {
	author := entity.Author{
		ID:   uuid.New().String(),
		Name: authorName,
	}

	createdAuthor, err := l.authorRepository.CreateAuthor(ctx, author)
	if err != nil {
		l.logger.Error("Failed to register author",
			zap.String("authorName", authorName),
			zap.Error(err))
		return entity.Author{}, err
	}

	return createdAuthor, nil
}

func (l *impl) GetAuthor(ctx context.Context, authorID string) (entity.Author, error) {
	author, err := l.authorRepository.GetAuthor(ctx, authorID)

	if err != nil {
		l.logger.Error("Failed to get author", zap.String("authorID", authorID), zap.Error(err))
		return entity.Author{}, err
	}
	return author, nil
}

func (l *impl) ChangeAuthor(ctx context.Context, authorID string, newAuthorName string) error {
	author := entity.Author{
		ID:   authorID,
		Name: newAuthorName,
	}

	if err := l.authorRepository.ChangeAuthor(ctx, author); err != nil {
		l.logger.Error("Failed to change author", zap.String("authorID", authorID), zap.Error(err))
		return err
	}

	return nil
}
