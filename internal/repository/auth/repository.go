package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

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

	//go:embed sql/get_by_hash.sql
	sqlGetByHash string

	//go:embed sql/delete_expired.sql
	sqlDeleteExpired string

	//go:embed sql/revoke.sql
	sqlRevoke string

	//go:embed sql/revoke_all_by_user.sql
	sqlRevokeAllByUser string
)

type repository struct {
	db *postgres.Postgres
}

func New(postgres *postgres.Postgres) Repository {
	return &repository{postgres}
}

func (r *repository) Create(ctx context.Context, in RefreshTokenIn) (id uuid.UUID, err error) {
	err = r.db.Pool.QueryRow(ctx, sqlCreate, in.UserID, in.TokenHash, time.Now().Add(in.TTL)).Scan(&id)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%w: failed to scan after insert", ErrDatabase)
		}
		return uuid.Nil, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}
	return id, nil
}

func (r *repository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.db.Pool.Exec(ctx, sqlDeleteExpired)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return 0, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return 0, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}
	return result.RowsAffected(), nil
}

func (r *repository) GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	token := models.RefreshToken{}
	err := r.db.Pool.QueryRow(ctx, sqlGetByHash, hash).Scan(
		&token.TokenID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
		&token.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: token with hash: %s not found", ErrTokenNotFound, hash)
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
	return &token, nil
}

func (r *repository) GetByID(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error) {
	token := models.RefreshToken{}
	err := r.db.Pool.QueryRow(ctx, sqlGetByID, tokenID).Scan(
		&token.TokenID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
		&token.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: token with id: %s not found", ErrTokenNotFound, tokenID)
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
	return &token, nil
}

func (r *repository) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	result, err := r.db.Pool.Exec(ctx, sqlRevoke, tokenID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: token with id: %s not found or already revoked", ErrCannotRevoke, tokenID)
	}

	return nil
}

func (r *repository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	result, err := r.db.Pool.Exec(ctx, sqlRevokeAllByUser, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return 0, fmt.Errorf(
				"%w: database error code %s: %v",
				ErrDatabase,
				pgErr.Code,
				pgErr.Message,
			)
		}
		return 0, fmt.Errorf("%w: query execution failed: %v", ErrDatabase, err)
	}

	return result.RowsAffected(), nil
}
