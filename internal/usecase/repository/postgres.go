package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/internal/entity"
)

var _ BooksRepository = (*postgresRepository)(nil)
var _ AuthorRepository = (*postgresRepository)(nil)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{db: db}
}
func (p *postgresRepository) CreateBook(ctx context.Context, book entity.Book) (entity.Book, error) {

	if err := p.validateAuthorsExist(ctx, book.AuthorIDs); err != nil {
		return entity.Book{}, err
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return entity.Book{}, err
	}
	defer tx.Rollback(ctx)

	const queryBook = `
		INSERT INTO book (id, name, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		RETURNING created_at, updated_at
	`
	result := entity.Book{
		ID:        book.ID,
		Name:      book.Name,
		AuthorIDs: book.AuthorIDs,
	}

	err = tx.QueryRow(ctx, queryBook, book.ID, book.Name).Scan(&result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return entity.Book{}, err
	}

	const queryAuthorBooks = `
		INSERT INTO author_book (author_id, book_id)
		VALUES ($1, $2)
	`

	for _, authorID := range book.AuthorIDs {
		_, err = tx.Exec(ctx, queryAuthorBooks, authorID, book.ID)
		if err != nil {
			return entity.Book{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return entity.Book{}, err
	}

	return result, nil
}

func (p *postgresRepository) UpdateBook(ctx context.Context, book entity.Book) error {
	if err := p.validateAuthorsExist(ctx, book.AuthorIDs); err != nil {
		return err
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const checkBookQuery = `
		SELECT 1
		FROM book
		WHERE id = $1
	`
	var exists int
	err = tx.QueryRow(ctx, checkBookQuery, book.ID).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.ErrBookNotFound
		}
		return err
	}

	const updateBookQuery = `
		UPDATE book
		SET name = $1, updated_at = now()
		WHERE id = $2
	`
	tag, err := tx.Exec(ctx, updateBookQuery, book.Name, book.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no book found with id %s", book.ID)
	}

	const deleteQuery = `
		DELETE FROM author_book
		WHERE book_id = $1
	`
	_, err = tx.Exec(ctx, deleteQuery, book.ID)
	if err != nil {
		return err
	}

	const insertQuery = `
		INSERT INTO author_book (author_id, book_id)
		VALUES ($1, $2)
	`
	for _, authorID := range book.AuthorIDs {
		_, err = tx.Exec(ctx, insertQuery, authorID, book.ID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) validateAuthorsExist(ctx context.Context, authorIDs []string) error {
	for _, authorID := range authorIDs {
		const query = `
			SELECT 1
			FROM author
			WHERE id = $1
		`
		var exists int
		err := p.db.QueryRow(ctx, query, authorID).Scan(&exists)
		if err != nil {
			if err == pgx.ErrNoRows {
				return entity.ErrAuthorNotFound
			}
			return err
		}
	}
	return nil
}

func (p *postgresRepository) GetBook(ctx context.Context, bookID string) (entity.Book, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM book
		WHERE id = $1
	`
	var book entity.Book
	err := p.db.QueryRow(ctx, query, bookID).Scan(&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Book{}, entity.ErrBookNotFound
	}
	if err != nil {
		return entity.Book{}, err
	}

	const queryAuthors = `
		SELECT author_id
		FROM author_book
		WHERE book_id = $1
	`
	rows, err := p.db.Query(ctx, queryAuthors, bookID)
	if err != nil {
		return entity.Book{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var authorID string
		if err := rows.Scan(&authorID); err != nil {
			return entity.Book{}, err
		}
		book.AuthorIDs = append(book.AuthorIDs, authorID)
	}
	return book, nil
}

func (p *postgresRepository) GetBooksByAuthor(ctx context.Context, authorID string) ([]entity.Book, error) {
	const query = `
		SELECT b.id, b.name, b.created_at, b.updated_at,
			(SELECT array_agg(ab2.author_id) FROM author_book ab2 WHERE ab2.book_id = b.id) AS author_ids
		FROM book b
		JOIN author_book ab ON ab.book_id = b.id
		WHERE ab.author_id = $1
		GROUP BY b.id, b.name, b.created_at, b.updated_at
	`
	rows, err := p.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		err := rows.Scan(&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.AuthorIDs)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (p *postgresRepository) CreateAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	const query = `
		INSERT INTO author (id, name)
		VALUES ($1, $2)
		RETURNING created_at, updated_at
	`
	err := p.db.QueryRow(ctx, query, author.ID, author.Name).
		Scan(&author.CreatedAt, &author.UpdatedAt)
	if err != nil {
		return entity.Author{}, err
	}
	return author, nil
}

func (p *postgresRepository) GetAuthor(ctx context.Context, authorID string) (entity.Author, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM author
		WHERE id = $1
	`
	var author entity.Author
	err := p.db.QueryRow(ctx, query, authorID).
		Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Author{}, entity.ErrAuthorNotFound
	}
	return author, err
}

func (p *postgresRepository) ChangeAuthor(ctx context.Context, author entity.Author) error {
	const query = `
		UPDATE author
		SET name = $1, updated_at = now()
		WHERE id = $2
	`
	tag, err := p.db.Exec(ctx, query, author.Name, author.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no author found with id %s", author.ID)
	}
	return nil
}
