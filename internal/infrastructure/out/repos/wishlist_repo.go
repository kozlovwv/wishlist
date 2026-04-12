package repos

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wishlist/internal/domain"
)

type WishlistRepo struct {
	db *pgxpool.Pool
	sq sq.StatementBuilderType
}

func NewWishlistRepo(db *pgxpool.Pool) *WishlistRepo {
	return &WishlistRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *WishlistRepo) Create(ctx context.Context, w domain.Wishlist) (domain.Wishlist, error) {
	query, args, err := r.sq.
		Insert("wishlists").
		Columns("user_id", "title", "description", "event_date", "public_token").
		Values(w.UserID, w.Title, w.Description, w.EventDate, w.PublicToken).
		Suffix("RETURNING id, user_id, title, description, event_date, public_token").
		ToSql()
	if err != nil {
		return domain.Wishlist{}, err
	}

	return r.scan(r.db.QueryRow(ctx, query, args...))
}

func (r *WishlistRepo) FindByID(ctx context.Context, id int64) (domain.Wishlist, error) {
	query, args, err := r.sq.
		Select("id", "user_id", "title", "description", "event_date", "public_token").
		From("wishlists").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, err
	}

	w, err := r.scan(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wishlist{}, domain.ErrWishlistNotFound
		}
		return domain.Wishlist{}, err
	}

	return w, nil
}

func (r *WishlistRepo) FindByUserID(ctx context.Context, userID int64) ([]domain.Wishlist, error) {
	query, args, err := r.sq.
		Select("id", "user_id", "title", "description", "event_date", "public_token").
		From("wishlists").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []domain.Wishlist
	for rows.Next() {
		var w domain.Wishlist
		if err := rows.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.PublicToken); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (r *WishlistRepo) FindByPublicToken(ctx context.Context, token string) (domain.Wishlist, error) {
	query, args, err := r.sq.
		Select("id", "user_id", "title", "description", "event_date", "public_token").
		From("wishlists").
		Where(sq.Eq{"public_token": token}).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, err
	}

	w, err := r.scan(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wishlist{}, domain.ErrWishlistNotFound
		}
		return domain.Wishlist{}, err
	}

	return w, nil
}

func (r *WishlistRepo) Update(ctx context.Context, w domain.Wishlist) (domain.Wishlist, error) {
	query, args, err := r.sq.
		Update("wishlists").
		Set("title", w.Title).
		Set("description", w.Description).
		Set("event_date", w.EventDate).
		Where(sq.Eq{"id": w.ID}).
		Suffix("RETURNING id, user_id, title, description, event_date, public_token").
		ToSql()
	if err != nil {
		return domain.Wishlist{}, err
	}

	return r.scan(r.db.QueryRow(ctx, query, args...))
}

func (r *WishlistRepo) Delete(ctx context.Context, id int64) error {
	query, args, err := r.sq.
		Delete("wishlists").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var deletedID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrWishlistNotFound
		}
		return err
	}

	return nil
}

func (r *WishlistRepo) scan(row pgx.Row) (domain.Wishlist, error) {
	var w domain.Wishlist
	err := row.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.PublicToken)
	return w, err
}
