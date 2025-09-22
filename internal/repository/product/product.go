package product

import (
	"context"
	"errors"
	"prodcrud/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateProduct(ctx context.Context, p *models.Product) error {
	_, err := r.db.Exec(ctx, `
	INSERT INTO products(name, price, quantity, description)
	VALUES ($1, $2, $3, $4)
	RETURNING id`, p.Name, p.Price, p.Quantity, p.Description)
	if err != nil {
		return errors.New("Failed to insert the product: " + err.Error())
	}
	return nil
}

func (r *Repo) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(ctx, `
	select id, name, price, quantity, description, created_at, updated_at from products where id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("failed to get product: " + err.Error())
	}
	return &p, nil
}

func (r *Repo) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT id, name, price, quantity, description, created_at, updated_at FROM products where deleted_at is null`)
	if err != nil {
		return nil, errors.New("failed to get products: " + err.Error())
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, errors.New("failed to scan products: " + err.Error())
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *Repo) UpdateProduct(ctx context.Context, p *models.Product) error {
	upd, err := r.db.Exec(ctx, `
	UPDATE products SET name = $1, price = $2, quantity = $3, description = $4,updated_at = now() 
	                WHERE id = $5 AND deleted_at is null`, p.Name, p.Price, p.Quantity, p.Description, p.ID)
	if err != nil {
		return errors.New("failed to update product: " + err.Error())
	}
	if upd.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *Repo) DeleteProduct(ctx context.Context, id int64) error {
	dlt, err := r.db.Exec(ctx, `
	UPDATE products SET deleted_at = now() WHERE id = $1
	`, id)
	if err != nil {
		return errors.New("failed to delete product: " + err.Error())
	}
	if dlt.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *Repo) RestoreProduct(ctx context.Context, id int64) error {
	restore, err := r.db.Exec(ctx, `
	UPDATE products SET deleted_at = null WHERE id = $1
	`, id)
	if err != nil {
		return errors.New("failed to restore product: " + err.Error())
	}
	if restore.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
