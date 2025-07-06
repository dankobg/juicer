package server

import (
	"context"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/core"
)

func (h *ApiHandler) ListGameTimeCategories(ctx context.Context, request api.ListGameTimeCategoriesRequestObject) (api.ListGameTimeCategoriesResponseObject, error) {
	gameTimeCategories, err := h.store.GameTimeCategory().List(ctx)
	if err != nil {
		return make(api.ListGameTimeCategories200JSONResponse, 0), nil
	}

	resp := make(api.ListGameTimeCategories200JSONResponse, 0)
	for _, x := range gameTimeCategories {
		resp = append(resp, api.GameTimeCategory{
			ID:                 x.ID,
			Name:               x.Name,
			UpperTimeLimitSecs: x.UpperTimeLimitSecs,
		})
	}

	return resp, nil
}

func (h *ApiHandler) ListGameTimeKinds(ctx context.Context, request api.ListGameTimeKindsRequestObject) (api.ListGameTimeKindsResponseObject, error) {
	gameTimeKinds, err := h.store.GameTimeKind().List(ctx)
	if err != nil {
		return make(api.ListGameTimeKinds200JSONResponse, 0), nil
	}

	resp := make(api.ListGameTimeKinds200JSONResponse, 0)
	for _, x := range gameTimeKinds {
		resp = append(resp, api.GameTimeKind{
			ID:      x.ID,
			Name:    x.Name,
			Enabled: x.Enabled,
		})
	}

	return resp, nil
}

func (h *ApiHandler) ListGameVariants(ctx context.Context, request api.ListGameVariantsRequestObject) (api.ListGameVariantsResponseObject, error) {
	gameVariants, err := h.store.GameVariant().List(ctx)
	if err != nil {
		return make(api.ListGameVariants200JSONResponse, 0), nil
	}

	resp := make(api.ListGameVariants200JSONResponse, 0)
	for _, x := range gameVariants {
		resp = append(resp, api.GameVariant{
			ID:      x.ID,
			Name:    x.Name,
			Enabled: x.Enabled,
		})
	}

	return resp, nil
}

func (h *ApiHandler) ListQuickGames(ctx context.Context, request api.ListQuickGamesRequestObject) (api.ListQuickGamesResponseObject, error) {
	return api.ListQuickGames200JSONResponse(core.QuickGames), nil
}

func (h *ApiHandler) GetGameStats(ctx context.Context, request api.GetGameStatsRequestObject) (api.GetGameStatsResponseObject, error) {
	gameStats, err := h.store.Game().GetGameStatsForUser(ctx, request.UserID, nil)
	if err != nil {
		return api.GetGameStats404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Message: "stats not found for user", Code: 404}}, nil
	}
	resp := api.GetGameStats200JSONResponse(gameStats)
	return resp, nil
}
