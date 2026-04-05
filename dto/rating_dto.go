package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func RatingToResponse(r models.Rating) api.Rating {
	return api.Rating{
		ID:                 r.ID,
		UserID:             r.UserID,
		GameTimeCategoryID: r.GameTimeCategoryID,
		Glicko:             int64(r.Glicko),
		Glicko2:            int64(r.Glicko2),
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}
