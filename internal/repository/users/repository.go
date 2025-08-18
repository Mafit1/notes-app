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
)

type repository struct {
	db *postgres.Postgres
}

func New(postgres *postgres.Postgres) Repository {
	return &repository{postgres}
}

// TODO: Сделать проверку на существование юзера
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

func (r *repository) GetByEmail(ctx context.Context, email string) (user models.User, err error) {
	row := r.db.Pool.QueryRow(ctx, sqlGetByEmail, email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%w: user with email: %s not found", ErrUserNotFound, email)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return models.User{}, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}

		return models.User{}, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	return user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (user models.User, err error) {
	row := r.db.Pool.QueryRow(ctx, sqlGetByID, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%w: user with id: %s not found", ErrUserNotFound, id)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return models.User{}, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}

		return models.User{}, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	return user, nil
}
