package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/store"
	"github.com/google/uuid"

	p "github.com/go-jet/jet/v2/postgres"
)

var _ store.UserStore = (*PgUserStore)(nil)

type PgUserStore struct {
	*PgStore
}

func NewPgUserStore(ps *PgStore) *PgUserStore {
	return &PgUserStore{
		PgStore: ps,
	}
}

func (s *PgUserStore) Create(ctx context.Context, in dto.UserChangeset) (model.User, error) {
	_, m := in.ToModel()

	q := t.User.INSERT(t.User.ID, t.User.MutableColumns.Except(t.User.CreatedAt, t.User.UpdatedAt)).
		MODEL(m).
		RETURNING(t.User.AllColumns)

	var dest model.User
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a user: %w", err)
	}
	return dest, nil
}

func (s *PgUserStore) Update(ctx context.Context, in dto.UserChangeset) (model.User, error) {
	cols, m := in.ToModel()

	q := t.User.UPDATE(cols).
		MODEL(m).
		RETURNING(t.User.AllColumns)

	var dest model.User
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a user: %w", err)
	}
	return dest, nil
}

func (s *PgUserStore) Delete(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	q := t.User.DELETE().WHERE(t.User.ID.EQ(p.UUID(userID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a user: %w", err)
	}
	return userID, nil
}

func (s *PgUserStore) Get(ctx context.Context, userID uuid.UUID) (model.User, error) {
	q := p.SELECT(t.User.AllColumns).
		FROM(t.User).
		WHERE(t.User.ID.EQ(p.UUID(userID)))

	var dest model.User
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a user: %w", err)
	}
	return dest, nil
}

func (s *PgUserStore) List(ctx context.Context) ([]model.User, error) {
	q := p.SELECT(t.User.AllColumns).
		FROM(t.User)

	var dest []model.User
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return dest, nil
}
