package server

import (
	"context"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/features/game"
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
	gameVariants, err := a.game.ListGameVariants(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list gamevariants: %w", err)
	}

	resp := api.ListGameVariants200JSONResponse{
		Data: gameVariants.Data,
		Meta: gameVariants.Meta.ToResp(),
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
	gameTimeCategories, err := a.game.ListGameTimeCategories(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list gametimecategories: %w", err)
	}

	resp := api.ListGameTimeCategories200JSONResponse{
		Data: gameTimeCategories.Data,
		Meta: gameTimeCategories.Meta.ToResp(),
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
	gameTimeKinds, err := a.game.ListGameTimeKinds(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list gametimekinds: %w", err)
	}

	resp := api.ListGameTimeKinds200JSONResponse{
		Data: gameTimeKinds.Data,
		Meta: gameTimeKinds.Meta.ToResp(),
	}

	return resp, nil
}

func (a *ApiHandler) GetGame(ctx context.Context, request api.GetGameRequestObject) (api.GetGameResponseObject, error) {
	game, err := a.game.GetGame(ctx, request)
	if err != nil {
		// if errors.Is(err, postgres.ErrGameNotFound) {
		// 	return api.GetGame404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("game_not_found", "game not found")}, nil
		// }
		return nil, fmt.Errorf("failed to get game by id: %w", err)
	}

	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	Tuple: &rts.RelationTuple{
	// 		Namespace: "Game",
	// 		Object:    shared.AuthzGameID(request.ID),
	// 		Relation:  "view",
	// 		Subject:   rts.NewSubjectID("*"),
	// 	},
	// }); err != nil || !checkResp.GetAllowed() {
	// 	return api.GetGamedefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "game_permission", "permission denied")}, nil
	// }
	resp := api.GetGame200JSONResponse(*game)

	return resp, nil
}

func (a *ApiHandler) ListGames(ctx context.Context, request api.ListGamesRequestObject) (api.ListGamesResponseObject, error) {
	// 	// if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
	// 	// 	Tuple: &rts.RelationTuple{
	// 	// 		Namespace: "Games",
	// 	// 		Object:    "games",
	// 	// 		Relation:  "view",
	// 	// 		Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
	// 	// 	},
	// 	// }); err != nil || !checkResp.GetAllowed() {
	// 	// 	return api.ListGames403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("game_permission", "permission denied")}, nil
	// 	// }
	games, err := a.game.ListGames(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	resp := api.ListGames200JSONResponse{
		Data: games.Data,
		Meta: games.Meta.ToResp(),
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
	ratings, err := a.game.ListRatings(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}

	resp := api.ListRatings200JSONResponse{
		Data: ratings.Data,
		Meta: ratings.Meta.ToResp(),
	}

	return resp, nil
}

func (a *ApiHandler) ListQuickGames(ctx context.Context, request api.ListQuickGamesRequestObject) (api.ListQuickGamesResponseObject, error) {
	return api.ListQuickGames200JSONResponse(game.QuickGames), nil
}

func (a *ApiHandler) GetGameStats(ctx context.Context, request api.GetGameStatsRequestObject) (api.GetGameStatsResponseObject, error) {
	gameStats, err := a.game.GetGameStats(ctx, request)
	if err != nil {
		return api.GetGameStatsdefaultJSONResponse{StatusCode: 500, Body: newGenericErr(500, "GAME_STATS", "failed to get stats")}, nil
	}
	return api.GetGameStats200JSONResponse(gameStats), nil
}
