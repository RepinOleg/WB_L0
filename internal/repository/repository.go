package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetOrderByID(id uint64) (string, error) {
	var (
		res string
		err error
	)
	row := r.db.QueryRow(`SELECT * FROM orders WHERE id = $1`, id)
	err = row.Scan(&id, &res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *Repository) AddOrder(order []byte) error {
	_, err := r.db.Exec(`INSERT INTO orders(order_data) VALUES($1)`, order)
	if err != nil {
		return err
	}
	return nil
}
