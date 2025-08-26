package users

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	//go:embed sql/create.sql
	sqlCreate string

	//go:embed sql/get_by_id.sql
	sqlGetByID string

	//go:embed sql/get_by_email.sql
	sqlGetByEmail string

	//go:embed sql/check_email_exists.sql
	sqlCheckEmailExists string
)

type repository struct {
	db *postgres.Postgres
}

func New(postgres *postgres.Postgres) Repository {
	return &repository{postgres}
}

func (r *repository) Create(ctx context.Context, user CreateUser) (id uuid.UUID, err error) {
	err = r.db.Pool.QueryRow(
		ctx,
		sqlCreate,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return uuid.Nil, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return uuid.Nil, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}
	return id, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}
	row := r.db.Pool.QueryRow(ctx, sqlGetByEmail, email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user with email: %s not found", ErrUserNotFound, email)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}

		return nil, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := models.User{}
	row := r.db.Pool.QueryRow(ctx, sqlGetByID, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user with id: %s not found", ErrUserNotFound, id)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}

		return nil, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	return &user, nil
}

func (r *repository) EmailExists(ctx context.Context, email string) (exists bool, err error) {
	err = r.db.Pool.QueryRow(ctx, sqlCheckEmailExists, email).Scan(&exists)
	return exists, err
}
