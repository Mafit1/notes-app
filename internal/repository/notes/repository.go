package notes

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/internal/repository"
	"github.com/Mafit1/notes-app/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	//go:embed sql/create.sql
	sqlCreate string

	//go:embed sql/get_all.sql
	sqlGetAll string

	//go:embed sql/get_by_id.sql
	sqlGetById string

	//go:embed sql/delete.sql
	sqlDelete string

	//go:embed sql/update.sql
	sqlUpdate string
)

type Repository struct {
	db *postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, note models.Note) (id int64, err error) {
	err = r.db.Pool.QueryRow(ctx, sqlCreate, note.Title, note.Content).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return 0, fmt.Errorf(
				"%w: database error code %s: %v",
				repository.ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: failed to scan after insert", repository.ErrDatabase)
		}
		return 0, fmt.Errorf("%w: query execution failed: %v", repository.ErrDatabase, err)
	}
	return id, nil
}

func (r *Repository) GetAll(ctx context.Context) (notes []models.Note, err error) {
	rows, err := r.db.Pool.Query(ctx, sqlGetAll)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf(
				"%w: database error code %s: %v",
				repository.ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return nil, fmt.Errorf("%w: query execution failed: %v", repository.ErrDatabase, err)
	}
	defer rows.Close()

	notes = make([]models.Note, 0)
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return nil, fmt.Errorf("%w: failed to scan row: %v", repository.ErrDatabase, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: rows iteration error: %v", repository.ErrDatabase, err)
	}

	if len(notes) == 0 {
		return []models.Note{}, nil
	}

	return notes, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (note *models.Note, err error) {
	row := r.db.Pool.QueryRow(ctx, sqlGetById, id)

	if err := row.Scan(&note); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf(
				"%w: database error code %s: %v",
				repository.ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: note with id: %d not found", repository.ErrNoteNotFound, id)
		}
		return nil, fmt.Errorf("%w: query execution failed: %v", repository.ErrDatabase, err)
	}

	return note, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Pool.Exec(ctx, sqlDelete, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf(
				"%w: database error code %s: %v",
				repository.ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return fmt.Errorf("%w: query execution failed: %v", repository.ErrDatabase, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: note with id: %d not found", repository.ErrNoteNotFound, id)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, note models.Note) error {
	result, err := r.db.Pool.Exec(
		ctx,
		sqlUpdate,
		note.ID,
		note.Title,
		note.Content,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf(
				"%w: database error code %s: %v",
				repository.ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return fmt.Errorf("%w: query execution failed: %v", repository.ErrDatabase, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: note with id: %d not found", repository.ErrNoteNotFound, note.ID)
	}

	return nil
}
