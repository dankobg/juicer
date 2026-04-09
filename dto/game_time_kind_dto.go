package dto

import (
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
)

func GameTimeKindToResponse(tk models.GameTimeKind) api.GameTimeKind {
	return api.GameTimeKind{
		ID:        tk.ID,
		Name:      tk.Name,
		Enabled:   tk.Enabled,
		CreatedAt: tk.CreatedAt,
		UpdatedAt: tk.UpdatedAt,
	}
}
