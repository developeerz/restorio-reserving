package postgres

import (
	"github.com/developeerz/restorio-reserving/internal/dto"
	"github.com/developeerz/restorio-reserving/internal/port"

	"github.com/jmoiron/sqlx"
)

// PostgresAdminRepository — реализация порта AdminRepository
type PostgresAdminRepository struct {
	db *sqlx.DB
}

// NewPostgresAdminRepository создаёт новый репозиторий
func NewPostgresAdminRepository(db *sqlx.DB) port.AdminRepository {
	return &PostgresAdminRepository{db: db}
}

// CreateTable сохраняет новый столик в таблице и возвращает его ID
func (r *PostgresAdminRepository) CreateTable(req dto.CreateTableRequest) (int, error) {
	var tableID int
	query := `
		INSERT INTO tables (restaurant_id, table_number, seats_number, type, shape)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING table_id
	`
	err := r.db.QueryRow(
		query,
		req.RestaurantID,
		req.TableNumber,
		req.SeatsNumber,
		req.Type,
		req.Shape,
	).Scan(&tableID)
	return tableID, err
}
