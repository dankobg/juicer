package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/store"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// var _ store.Store = (*PgStore)(nil)

type PgStore struct {
	db *sql.DB
}

// Connect creates the new connection
func Connect(dbSettings config.DatabaseConfig) (*sql.DB, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbSettings.User, dbSettings.Password),
		Host:     net.JoinHostPort(dbSettings.Host, strconv.Itoa(dbSettings.Port)),
		Path:     dbSettings.DB,
		RawQuery: url.Values{"sslmode": []string{dbSettings.SSLMode}}.Encode(),
	}

	var db *sql.DB
	var err error

	db, err = sql.Open("pgx", dsn.String())
	for err != nil {
		log.Printf("failed to connect to database")

		if dbSettings.RetriesNum > 0 {
			dbSettings.RetriesNum--
			log.Printf("retrying the database connection. retries left: (%d)", dbSettings.RetriesNum)
			time.Sleep(dbSettings.RetriesDelay)
			db, err = sql.Open("pgx", dsn.String())
			continue
		}

		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}

// New creates the new store
func New(db *sql.DB) *PgStore {
	return &PgStore{db: db}
}

func (ps *PgStore) User() store.UserStore {
	return NewPgUserStore(ps)
}

func (ps *PgStore) Game() store.GameStore {
	return NewPgGameStore(ps)
}

func (ps *PgStore) GameResult() store.GameResultStore {
	return NewPgGameResultStore(ps)
}

func (ps *PgStore) GameResultStatus() store.GameResultStatusStore {
	return NewPgGameResultStatusStore(ps)
}

func (ps *PgStore) GameState() store.GameStateStore {
	return NewPgGameStateStore(ps)
}

func (ps *PgStore) GameTimeCategory() store.GameTimeCategoryStore {
	return NewPgGameTimeCategoryStore(ps)
}

func (ps *PgStore) GameTimeKind() store.GameTimeKindStore {
	return NewPgGameTimeKindStore(ps)
}

func (ps *PgStore) GameVariant() store.GameVariantStore {
	return NewPgGameVariantStore(ps)
}

func (ps *PgStore) Rating() store.RatingStore {
	return NewPgRatingStore(ps)
}

func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}
