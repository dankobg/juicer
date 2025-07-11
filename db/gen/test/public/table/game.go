//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Game = newGameTable("public", "game", "")

type gameTable struct {
	postgres.Table

	// Columns
	ID                   postgres.ColumnString
	WhiteID              postgres.ColumnString
	BlackID              postgres.ColumnString
	GuestWhiteID         postgres.ColumnString
	GuestBlackID         postgres.ColumnString
	VariantID            postgres.ColumnString
	TimeKindID           postgres.ColumnString
	TimeCategoryID       postgres.ColumnString
	IsGuest              postgres.ColumnBool
	TimeControlClock     postgres.ColumnInteger
	TimeControlIncrement postgres.ColumnInteger
	ReconnectTimeout     postgres.ColumnInteger
	FirstMoveTimeout     postgres.ColumnInteger
	WhiteGameClock       postgres.ColumnInteger
	BlackGameClock       postgres.ColumnInteger
	ResultID             postgres.ColumnString
	ResultStatusID       postgres.ColumnString
	StateID              postgres.ColumnString
	StartTime            postgres.ColumnTimestampz
	EndTime              postgres.ColumnTimestampz
	LastMove             postgres.ColumnTimestampz
	Fen                  postgres.ColumnString
	Pgn                  postgres.ColumnString
	CreatedAt            postgres.ColumnTimestampz
	UpdatedAt            postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
	DefaultColumns postgres.ColumnList
}

type GameTable struct {
	gameTable

	EXCLUDED gameTable
}

// AS creates new GameTable with assigned alias
func (a GameTable) AS(alias string) *GameTable {
	return newGameTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new GameTable with assigned schema name
func (a GameTable) FromSchema(schemaName string) *GameTable {
	return newGameTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new GameTable with assigned table prefix
func (a GameTable) WithPrefix(prefix string) *GameTable {
	return newGameTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new GameTable with assigned table suffix
func (a GameTable) WithSuffix(suffix string) *GameTable {
	return newGameTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newGameTable(schemaName, tableName, alias string) *GameTable {
	return &GameTable{
		gameTable: newGameTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newGameTableImpl("", "excluded", ""),
	}
}

func newGameTableImpl(schemaName, tableName, alias string) gameTable {
	var (
		IDColumn                   = postgres.StringColumn("id")
		WhiteIDColumn              = postgres.StringColumn("white_id")
		BlackIDColumn              = postgres.StringColumn("black_id")
		GuestWhiteIDColumn         = postgres.StringColumn("guest_white_id")
		GuestBlackIDColumn         = postgres.StringColumn("guest_black_id")
		VariantIDColumn            = postgres.StringColumn("variant_id")
		TimeKindIDColumn           = postgres.StringColumn("time_kind_id")
		TimeCategoryIDColumn       = postgres.StringColumn("time_category_id")
		IsGuestColumn              = postgres.BoolColumn("is_guest")
		TimeControlClockColumn     = postgres.IntegerColumn("time_control_clock")
		TimeControlIncrementColumn = postgres.IntegerColumn("time_control_increment")
		ReconnectTimeoutColumn     = postgres.IntegerColumn("reconnect_timeout")
		FirstMoveTimeoutColumn     = postgres.IntegerColumn("first_move_timeout")
		WhiteGameClockColumn       = postgres.IntegerColumn("white_game_clock")
		BlackGameClockColumn       = postgres.IntegerColumn("black_game_clock")
		ResultIDColumn             = postgres.StringColumn("result_id")
		ResultStatusIDColumn       = postgres.StringColumn("result_status_id")
		StateIDColumn              = postgres.StringColumn("state_id")
		StartTimeColumn            = postgres.TimestampzColumn("start_time")
		EndTimeColumn              = postgres.TimestampzColumn("end_time")
		LastMoveColumn             = postgres.TimestampzColumn("last_move")
		FenColumn                  = postgres.StringColumn("fen")
		PgnColumn                  = postgres.StringColumn("pgn")
		CreatedAtColumn            = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn            = postgres.TimestampzColumn("updated_at")
		allColumns                 = postgres.ColumnList{IDColumn, WhiteIDColumn, BlackIDColumn, GuestWhiteIDColumn, GuestBlackIDColumn, VariantIDColumn, TimeKindIDColumn, TimeCategoryIDColumn, IsGuestColumn, TimeControlClockColumn, TimeControlIncrementColumn, ReconnectTimeoutColumn, FirstMoveTimeoutColumn, WhiteGameClockColumn, BlackGameClockColumn, ResultIDColumn, ResultStatusIDColumn, StateIDColumn, StartTimeColumn, EndTimeColumn, LastMoveColumn, FenColumn, PgnColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns             = postgres.ColumnList{WhiteIDColumn, BlackIDColumn, GuestWhiteIDColumn, GuestBlackIDColumn, VariantIDColumn, TimeKindIDColumn, TimeCategoryIDColumn, IsGuestColumn, TimeControlClockColumn, TimeControlIncrementColumn, ReconnectTimeoutColumn, FirstMoveTimeoutColumn, WhiteGameClockColumn, BlackGameClockColumn, ResultIDColumn, ResultStatusIDColumn, StateIDColumn, StartTimeColumn, EndTimeColumn, LastMoveColumn, FenColumn, PgnColumn, CreatedAtColumn, UpdatedAtColumn}
		defaultColumns             = postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return gameTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                   IDColumn,
		WhiteID:              WhiteIDColumn,
		BlackID:              BlackIDColumn,
		GuestWhiteID:         GuestWhiteIDColumn,
		GuestBlackID:         GuestBlackIDColumn,
		VariantID:            VariantIDColumn,
		TimeKindID:           TimeKindIDColumn,
		TimeCategoryID:       TimeCategoryIDColumn,
		IsGuest:              IsGuestColumn,
		TimeControlClock:     TimeControlClockColumn,
		TimeControlIncrement: TimeControlIncrementColumn,
		ReconnectTimeout:     ReconnectTimeoutColumn,
		FirstMoveTimeout:     FirstMoveTimeoutColumn,
		WhiteGameClock:       WhiteGameClockColumn,
		BlackGameClock:       BlackGameClockColumn,
		ResultID:             ResultIDColumn,
		ResultStatusID:       ResultStatusIDColumn,
		StateID:              StateIDColumn,
		StartTime:            StartTimeColumn,
		EndTime:              EndTimeColumn,
		LastMove:             LastMoveColumn,
		Fen:                  FenColumn,
		Pgn:                  PgnColumn,
		CreatedAt:            CreatedAtColumn,
		UpdatedAt:            UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
