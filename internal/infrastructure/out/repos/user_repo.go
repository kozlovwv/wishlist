package repos

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wishlist/internal/domain"
)

type UserRepo struct {
	db *pgxpool.Pool
	sq sq.StatementBuilderType
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query, args, err := r.sq.
		Insert("users").
		Columns("email", "password_hash").
		Values(user.Email, user.PasswordHash).
		Suffix("RETURNING id, email, password_hash").
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("build create user query: %w", err)
	}

	var u domain.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("exec create user: %w", err)
	}

	return u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query, args, err := r.sq.
		Select("id", "email", "password_hash").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("build find user by email query: %w", err)
	}

	var u domain.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("query find user by email: %w", err)
	}

	return u, nil
}
