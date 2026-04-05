package server

import (
	"context"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/dto"
	"github.com/dankobg/juicer/persistence/dbtype"
)

func (a *ApiHandler) ListGameVariants(ctx context.Context, request api.ListGameVariantsRequestObject) (api.ListGameVariantsResponseObject, error) {
	// sess := GetSession(ctx)

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "GameVariants",
	// 		Object:    "gamevariants",
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.ListGameVariants403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("gamevariant_permission", "permission denied")}, nil
	// }
	filters := dbtype.ListGameVariantsFilters{ListGameVariantsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	gamevariants, err := a.persistor.GameVariant().ListGameVariants(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list gamevariants: %w", err)
	}

	gamevariantsData := make([]api.GameVariant, len(gamevariants.Data))
	for i, gamevariant := range gamevariants.Data {
		gamevariantsData[i] = dto.GameVariantToResponse(gamevariant)
	}

	resp := api.ListGameVariants200JSONResponse{
		Data: gamevariantsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, gamevariants.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListGameTimeCategories(ctx context.Context, request api.ListGameTimeCategoriesRequestObject) (api.ListGameTimeCategoriesResponseObject, error) {
	// sess := GetSession(ctx)

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "GameTimeCategories",
	// 		Object:    "gametimecategories",
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.ListGameTimeCategories403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("gametimecategory_permission", "permission denied")}, nil
	// }
	filters := dbtype.ListGameTimeCategoriesFilters{ListGameTimeCategoriesParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	gametimecategories, err := a.persistor.GameTimeCategory().ListGameTimeCategories(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list gametimecategories: %w", err)
	}

	gametimecategoriesData := make([]api.GameTimeCategory, len(gametimecategories.Data))
	for i, gametimecategory := range gametimecategories.Data {
		gametimecategoriesData[i] = dto.GameTimeCategoryToResponse(gametimecategory)
	}

	resp := api.ListGameTimeCategories200JSONResponse{
		Data: gametimecategoriesData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, gametimecategories.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListGameTimeKinds(ctx context.Context, request api.ListGameTimeKindsRequestObject) (api.ListGameTimeKindsResponseObject, error) {
	// sess := GetSession(ctx)

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "GameTimeKinds",
	// 		Object:    "gametimekinds",
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.ListGameTimeKinds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("gametimekind_permission", "permission denied")}, nil
	// }
	filters := dbtype.ListGameTimeKindsFilters{ListGameTimeKindsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	gametimekinds, err := a.persistor.GameTimeKind().ListGameTimeKinds(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list gametimekinds: %w", err)
	}

	gametimekindsData := make([]api.GameTimeKind, len(gametimekinds.Data))
	for i, gametimekind := range gametimekinds.Data {
		gametimekindsData[i] = dto.GameTimeKindToResponse(gametimekind)
	}

	resp := api.ListGameTimeKinds200JSONResponse{
		Data: gametimekindsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, gametimekinds.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListGames(ctx context.Context, request api.ListGamesRequestObject) (api.ListGamesResponseObject, error) {
	// sess := GetSession(ctx)

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "Games",
	// 		Object:    "games",
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.ListGames403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("game_permission", "permission denied")}, nil
	// }
	filters := dbtype.ListGamesFilters{ListGamesParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	games, err := a.persistor.Game().ListGames(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	gamesData := make([]api.Game, len(games.Data))
	for i, game := range games.Data {
		gamesData[i] = dto.GameToResponse(game)
	}

	resp := api.ListGames200JSONResponse{
		Data: gamesData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, games.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListRatings(ctx context.Context, request api.ListRatingsRequestObject) (api.ListRatingsResponseObject, error) {
	// sess := GetSession(ctx)

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "Ratings",
	// 		Object:    "ratings",
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.ListRatings403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("rating_permission", "permission denied")}, nil
	// }
	filters := dbtype.ListRatingsFilters{ListRatingsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	ratings, err := a.persistor.Rating().ListRatings(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}

	ratingsData := make([]api.Rating, len(ratings.Data))
	for i, rating := range ratings.Data {
		ratingsData[i] = dto.RatingToResponse(rating)
	}

	resp := api.ListRatings200JSONResponse{
		Data: ratingsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, ratings.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListQuickGames(ctx context.Context, request api.ListQuickGamesRequestObject) (api.ListQuickGamesResponseObject, error) {
	quickGames := []api.QuickGame{
		{Name: "Hyperbullet", ClockSecs: 30, IncrementSecs: 0},
		{Name: "Bullet", ClockSecs: 60, IncrementSecs: 0},
		{Name: "Blitz", ClockSecs: 180, IncrementSecs: 0},
		{Name: "Blitz", ClockSecs: 180, IncrementSecs: 1},
		{Name: "Blitz", ClockSecs: 300, IncrementSecs: 0},
		{Name: "Blitz", ClockSecs: 300, IncrementSecs: 2},
		{Name: "Rapid", ClockSecs: 600, IncrementSecs: 0},
		{Name: "Rapid", ClockSecs: 600, IncrementSecs: 5},
		{Name: "Rapid", ClockSecs: 900, IncrementSecs: 0},
		{Name: "Rapid", ClockSecs: 900, IncrementSecs: 5},
		{Name: "Classical", ClockSecs: 1800, IncrementSecs: 0},
		{Name: "Classical", ClockSecs: 2700, IncrementSecs: 10},
	}

	return api.ListQuickGames200JSONResponse(quickGames), nil
}

func (a *ApiHandler) GetGameStats(ctx context.Context, request api.GetGameStatsRequestObject) (api.GetGameStatsResponseObject, error) {
	return api.GetGameStats200JSONResponse{}, nil
	// 	gameStats, err := a.persistor.Game().GetGameStatsForUser(ctx, request.UserID, nil)
	// 	if err != nil {
	// 		return api.GetGameStats404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Message: "stats not found for user", Code: 404}}, nil
	// 	}
	// 	resp := api.GetGameStats200JSONResponse(gameStats)
	// 	return resp, nil
}
