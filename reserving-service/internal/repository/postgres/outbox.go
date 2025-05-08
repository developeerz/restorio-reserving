package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OutboxRepository struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
	return &OutboxRepository{
		db: db,
	}
}

func (r *OutboxRepository) UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE outbox
		SET send_status = true
		WHERE id = $1;
	`
	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
