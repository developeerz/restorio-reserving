package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OutboxRepository interface {
	UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error
	Transaction(ctx context.Context, fn func(repo OutboxRepository) error) error
}

type outboxRepository struct {
	db sqlx.ExtContext
}

func NewOutboxRepository(db *sqlx.DB) OutboxRepository {
	return &outboxRepository{
		db: db,
	}
}

type txoutboxRepository struct {
	tx *sqlx.Tx
}

func (r *outboxRepository) UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE outbox
		SET send_status = true
		WHERE id = $1;
	`
	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *outboxRepository) Transaction(ctx context.Context, fn func(repo OutboxRepository) error) error {
	db, ok := r.db.(*sqlx.DB)
	if !ok {
		return errors.New("cannot start transaction")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	var rbErr error

	defer func() {
		if rbErr != nil {
			if err := tx.Rollback(); err != nil {
				rbErr = fmt.Errorf("original error: %v, rollback error: %w", rbErr, err)
			}
		}
	}()

	repo := r.withTx(tx)

	err = fn(repo)
	if err != nil {
		rbErr = err
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (r *outboxRepository) withTx(tx *sqlx.Tx) *outboxRepository {
	return &outboxRepository{
		db: tx,
	}
}
