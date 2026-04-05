package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func GameTimeCategoryToResponse(tc models.GameTimeCategory) api.GameTimeCategory {
	return api.GameTimeCategory{
		ID:        tc.ID,
		Name:      tc.Name,
		CreatedAt: tc.CreatedAt,
		UpdatedAt: tc.UpdatedAt,
	}
}
