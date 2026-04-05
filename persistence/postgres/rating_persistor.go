package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ persistence.RatingPersistor = (*PgRatingPersistor)(nil)

type PgRatingPersistor struct {
	*PgPersistor
}

func NewPgRatingPersistor(ps *PgPersistor) *PgRatingPersistor {
	return &PgRatingPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrRatingIntegrityViolation struct{ errIntegrityViolation }
	ErrRatingUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrRatingNotFound         = errors.New("rating not found")
	errRatingIntegrity        = ErrRatingIntegrityViolation{}
	errRatingUniqueName       = ErrRatingUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errRatingUniqueIsoAlpha2  = ErrRatingUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errRatingUniqueIsoAlpha3  = ErrRatingUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errRatingUniqueIsoNumeric = ErrRatingUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertRatingPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_rating_name":
			return errRatingUniqueName
		case "uq_rating_iso_alpha2":
			return errRatingUniqueIsoAlpha2
		case "uq_rating_iso_alpha3":
			return errRatingUniqueIsoAlpha3
		case "uq_rating_iso_numeric":
			return errRatingUniqueIsoNumeric
		}

		return errRatingIntegrity
	}

	return pgErr
}

func (pst *PgRatingPersistor) ListRatings(ctx context.Context, filters dbtype.ListRatingsFilters) (dbtype.PagedResult[models.Rating], error) {
	q := psql.Select(
		sm.Columns(models.Ratings.Columns),
		sm.From(models.Ratings.Name()),
		sm.GroupBy(models.Ratings.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Ratings.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListRatingsParams) {
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

	ratings, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListRatingsRow]())
	if err != nil {
		return dbtype.PagedResult[models.Rating]{}, fmt.Errorf("query ratings")
	}

	result := dbtype.PagedResult[models.Rating]{
		Data: make([]models.Rating, len(ratings)),
	}
	for i, row := range ratings {
		result.Data[i] = row.Rating
	}

	if len(ratings) > 0 {
		result.TotalCount = ratings[0].TotalCount
	}

	return result, nil
}

func (pst *PgRatingPersistor) GetRatingByID(ctx context.Context, ratingID int64) (models.Rating, error) {
	q := psql.Select(
		sm.Columns(models.Ratings.Columns),
		sm.From(models.Ratings.Name()),
		sm.Where(models.Ratings.Columns.ID.EQ(psql.Arg(ratingID))),
	)

	rating, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return models.Rating{}, fmt.Errorf("query rating")
	}

	return rating, nil
}

func (pst *PgRatingPersistor) CreateRating(ctx context.Context, in models.RatingSetter) (models.Rating, error) {
	q := models.Ratings.Insert(&in, im.Returning(models.Ratings.Columns))

	rating, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Rating]())
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

	rating, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Rating]())
	if err != nil {
		return models.Rating{}, fmt.Errorf("update rating")
	}

	return rating, nil
}

func (pst *PgRatingPersistor) DeleteRatingByID(ctx context.Context, ratingID int64) (int64, error) {
	q := models.Ratings.Delete(dm.Where(models.Ratings.Columns.ID.EQ(psql.Arg(ratingID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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

	ratings, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.Rating]())
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete ratings: %w", err)
	}

	return nil
}
