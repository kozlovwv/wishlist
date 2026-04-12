package repos

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wishlist/internal/domain"
)

type ItemRepo struct {
	db *pgxpool.Pool
	sq sq.StatementBuilderType
}

func NewItemRepo(db *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ItemRepo) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	query, args, err := r.sq.
		Insert("items").
		Columns("wishlist_id", "title", "description", "url", "priority").
		Values(item.WishlistID, item.Title, item.Description, item.URL, item.Priority).
		Suffix("RETURNING id, wishlist_id, title, description, url, priority, is_reserved").
		ToSql()
	if err != nil {
		return domain.Item{}, err
	}

	return r.scan(r.db.QueryRow(ctx, query, args...))
}

func (r *ItemRepo) FindByID(ctx context.Context, id int64) (domain.Item, error) {
	query, args, err := r.sq.
		Select("id", "wishlist_id", "title", "description", "url", "priority", "is_reserved").
		From("items").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.Item{}, err
	}

	item, err := r.scan(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Item{}, domain.ErrItemNotFound
		}
		return domain.Item{}, err
	}

	return item, nil
}

func (r *ItemRepo) FindByWishlistID(ctx context.Context, wishlistID int64) ([]domain.Item, error) {
	query, args, err := r.sq.
		Select("id", "wishlist_id", "title", "description", "url", "priority", "is_reserved").
		From("items").
		Where(sq.Eq{"wishlist_id": wishlistID}).
		OrderBy("priority DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(
			&item.ID, &item.WishlistID, &item.Title, &item.Description,
			&item.URL, &item.Priority, &item.IsReserved,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepo) Update(ctx context.Context, item domain.Item) (domain.Item, error) {
	query, args, err := r.sq.
		Update("items").
		Set("title", item.Title).
		Set("description", item.Description).
		Set("url", item.URL).
		Set("priority", item.Priority).
		Where(sq.Eq{"id": item.ID}).
		Suffix("RETURNING id, wishlist_id, title, description, url, priority, is_reserved").
		ToSql()
	if err != nil {
		return domain.Item{}, err
	}

	return r.scan(r.db.QueryRow(ctx, query, args...))
}

func (r *ItemRepo) Delete(ctx context.Context, id int64) error {
	query, args, err := r.sq.
		Delete("items").
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
			return domain.ErrItemNotFound
		}
		return err
	}

	return nil
}

func (r *ItemRepo) Reserve(ctx context.Context, id int64) error {
	query, args, err := r.sq.
		Update("items").
		Set("is_reserved", true).
		Where(sq.Eq{"id": id, "is_reserved": false}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	var updatedID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrItemReserved
		}
		return err
	}

	return nil
}

func (r *ItemRepo) scan(row pgx.Row) (domain.Item, error) {
	var item domain.Item
	err := row.Scan(
		&item.ID, &item.WishlistID, &item.Title, &item.Description,
		&item.URL, &item.Priority, &item.IsReserved,
	)
	return item, err
}
