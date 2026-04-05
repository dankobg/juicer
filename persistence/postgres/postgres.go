package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/persistence"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
)

// var _ persistence.Persistor = (*PgPersistor)(nil)

type PgPersistor struct {
	pool *pgxpool.Pool
	exec bobpgx.Pool
}

func NewPool(ctx context.Context, dbSettings config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbSettings.User, dbSettings.Password),
		Host:     net.JoinHostPort(dbSettings.Host, strconv.Itoa(dbSettings.Port)),
		Path:     dbSettings.DB,
		RawQuery: url.Values{"sslmode": []string{dbSettings.SSLMode}}.Encode(),
	}

	poolCfg, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	return pool, nil
}

// New creates the new persistence
func New(pool *pgxpool.Pool) *PgPersistor {
	return &PgPersistor{pool: pool, exec: bobpgx.NewPool(pool)}
}

func (ps *PgPersistor) User() persistence.UserPersistor {
	return NewPgUserPersistor(ps)
}

func (ps *PgPersistor) Game() persistence.GamePersistor {
	return NewPgGamePersistor(ps)
}

func (ps *PgPersistor) GameResult() persistence.GameResultPersistor {
	return NewPgGameResultPersistor(ps)
}

func (ps *PgPersistor) GameResultStatus() persistence.GameResultStatusPersistor {
	return NewPgGameResultStatusPersistor(ps)
}

func (ps *PgPersistor) GameState() persistence.GameStatePersistor {
	return NewPgGameStatePersistor(ps)
}

func (ps *PgPersistor) GameTimeCategory() persistence.GameTimeCategoryPersistor {
	return NewPgGameTimeCategoryPersistor(ps)
}

func (ps *PgPersistor) GameTimeKind() persistence.GameTimeKindPersistor {
	return NewPgGameTimeKindPersistor(ps)
}

func (ps *PgPersistor) GameVariant() persistence.GameVariantPersistor {
	return NewPgGameVariantPersistor(ps)
}

func (ps *PgPersistor) Rating() persistence.RatingPersistor {
	return NewPgRatingPersistor(ps)
}

func (ps *PgPersistor) WithTx(
	ctx context.Context,
	fn func(tx bobpgx.Tx) error,
) (err error) {
	tx, err := ps.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	bobTx := bobpgx.NewTx(tx, func() {})

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)

			panic(r)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	return fn(bobTx)
}

func valOrNil[T any](v omitnull.Val[T]) *T {
	var out *T
	if !v.IsUnset() && !v.IsNull() {
		out = v.MustPtr()
	}

	return out
}
