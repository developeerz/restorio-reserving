package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/developeerz/restorio-reserving/internal/port"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type outboxRepo struct {
	db sqlx.ExtContext
}

func NewOutboxRepo(db *sqlx.DB) port.OutboxRepository {
	return &outboxRepo{db: db}
}

func (r *outboxRepo) UpdateSendStatusTrueByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE outbox SET send_status = true WHERE id = $1
    `, id)
	return err
}

func (r *outboxRepo) DeleteSent(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
        DELETE FROM outbox WHERE send_status = true
    `)
	return err
}

func (r *outboxRepo) Transaction(ctx context.Context, fn func(repo port.OutboxRepository) error) error {
	db, ok := r.db.(*sqlx.DB)
	if !ok {
		return errors.New("cannot start transaction")
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			panic(e)
		}
	}()
	repoTx := &outboxRepo{db: tx}
	if err := fn(repoTx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *outboxRepo) CreateOutbox(ctx context.Context, id, topic string, payload []byte, sendTime time.Time) error {
	q := `INSERT INTO outbox (id, topic, payload, send_time, send_status)
          VALUES ($1,$2,$3,$4,'pending')`
	_, err := r.db.ExecContext(ctx, q, id, topic, payload, sendTime)
	return err
}

func (r *outboxRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Outbox, error) {
	var e entity.Outbox
	err := sqlx.GetContext(ctx, r.db, &e, `
		SELECT id, topic, payload, send_time
		FROM outbox
		WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *outboxRepo) GetTablePayload(ctx context.Context, tableID int) (string, []byte, error) {
	var e entity.OutboxPayload
	q := `SELECT r.name, r.address, t.table_number
          FROM restaurants r
          JOIN tables t ON t.restaurant_id=r.restaurant_id
          WHERE t.table_id=$1`
	if err := sqlx.GetContext(ctx, r.db, &e, q, tableID); err != nil {
		return "", nil, err
	}
	b, err := entity.PayloadToJSON(e)
	return "reservation.created", b, err
}
