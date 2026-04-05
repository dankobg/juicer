package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func GameVariantToResponse(gv models.GameVariant) api.GameVariant {
	return api.GameVariant{
		ID:        gv.ID,
		Enabled:   gv.Enabled,
		Name:      gv.Name,
		CreatedAt: gv.CreatedAt,
		UpdatedAt: gv.UpdatedAt,
	}
}
