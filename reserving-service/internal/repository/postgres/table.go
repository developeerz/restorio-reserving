package postgres

import (
	"context"

	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/table"
	"github.com/jmoiron/sqlx"
)

type TableRepository struct {
	db sqlx.ExtContext
}

func NewTableRepository(db *sqlx.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) GetTablesByRestaurantID(ctx context.Context, restaurantID int) ([]table.Table, error) {
	query := `
		SELECT
			t.table_id,
			t.restaurant_id,
			t.table_number,
			t.seats_number,
			t.type,
			t.shape,
			p.x,
			p.y
		FROM tables t
		JOIN positions p ON p.table_id = t.table_id
		WHERE t.restaurant_id = $1;
	`

	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}

	tables := make([]table.Table, 0)

	defer rows.Close()

	for rows.Next() {
		var table table.Table

		if err = rows.Scan(
			&table.TableID,
			&table.TableNumber,
			&table.Type,
			&table.Shape,
			&table.X,
			&table.Y,
		); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}
