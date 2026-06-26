package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/pagination"
	"github.com/dankobg/juicer/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ game.RatingPersistor = (*PgRatingPersistor)(nil)

type PgRatingPersistor struct {
	*postgres.PgPersistor
}

func NewPgRatingPersistor(pst *postgres.PgPersistor) *PgRatingPersistor {
	return &PgRatingPersistor{
		PgPersistor: pst,
	}
}

func convertGameRatingPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgRatingPersistor) ListRatings(ctx context.Context, filters game.ListRatingsFilters) (pagination.WithTotal[models.Rating], error) {
	q := psql.Select(
		sm.Columns(models.Ratings.Columns),
		sm.From(models.Ratings.Name()),
		sm.GroupBy(models.Ratings.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.Ratings.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters.ListRatingsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Ratings.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListRatingsRow struct {
		models.Rating
		TotalCount int64
	}

	ratingsRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListRatingsRow]())
	if err != nil {
		return pagination.WithTotal[models.Rating]{}, fmt.Errorf("query ratings")
	}

	ratings := make([]models.Rating, len(ratingsRows))
	for i, row := range ratingsRows {
		ratings[i] = row.Rating
	}

	var total int64
	if len(ratingsRows) > 0 {
		total = ratingsRows[0].TotalCount
	}

	out := pagination.NewWithTotal(ratings, total)

	return out, nil
}

func (pst *PgRatingPersistor) GetRatingByID(ctx context.Context, ratingID int64) (models.Rating, error) {
	q := psql.Select(
		sm.Columns(models.Ratings.Columns),
		sm.From(models.Ratings.Name()),
		sm.Where(models.Ratings.Columns.ID.EQ(psql.Arg(ratingID))),
	)

	rating, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return models.Rating{}, fmt.Errorf("query rating")
	}

	return rating, nil
}

func (pst *PgRatingPersistor) CreateRating(ctx context.Context, in models.RatingSetter) (models.Rating, error) {
	q := models.Ratings.Insert(&in, im.Returning(models.Ratings.Columns))

	rating, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return models.Rating{}, fmt.Errorf("insert rating")
	}

	return rating, nil
}

func (pst *PgRatingPersistor) UpdateRating(ctx context.Context, ratingID int64, in models.RatingSetter) (models.Rating, error) {
	q := models.Ratings.Update(
		in.UpdateMod(),
		um.Where(models.Ratings.Columns.ID.EQ(psql.Arg(ratingID))),
		um.Returning(models.Ratings.Columns),
	)

	rating, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return models.Rating{}, fmt.Errorf("update rating")
	}

	return rating, nil
}

func (pst *PgRatingPersistor) DeleteRatingByID(ctx context.Context, ratingID int64) (int64, error) {
	q := models.Ratings.Delete(dm.Where(models.Ratings.Columns.ID.EQ(psql.Arg(ratingID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete rating: %w", err)
	}

	return ratingID, nil
}

func (pst *PgRatingPersistor) BulkCreateRatings(ctx context.Context, in []models.RatingSetter) ([]models.Rating, error) {
	setters := make([]*models.RatingSetter, len(in))
	for i, x := range in {
		setters[i] = &x
	}

	q := models.Ratings.Insert(bob.ToMods(setters...), im.Returning(models.Ratings.Columns))

	ratings, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return nil, fmt.Errorf("bulk insert ratings")
	}

	return ratings, nil
}

func (pst *PgRatingPersistor) BulkDeleteRatings(ctx context.Context, ids []int64) error {
	ratingIDs := make([]any, len(ids))
	for i, id := range ids {
		ratingIDs[i] = id
	}

	q := models.Ratings.Delete(dm.Where(models.Ratings.Columns.ID.In(psql.Arg(ratingIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete ratings: %w", err)
	}

	return nil
}
