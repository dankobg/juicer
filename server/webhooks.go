package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/opt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (h *ApiHandler) registrationAfterPassword(c echo.Context) error {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := c.Bind(&payload); err != nil {
		h.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_password"), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	h.Log.Debug("kratos webhook registration_after_password", slog.String("identity_id", payload.Identity.Id))
	v := c.Request().Header.Get("Authorization")
	if v != h.Cfg.KratosAPIKey {
		h.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_password"), slog.String("authorization", v))
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err := h.createUserRelationTuples(c.Request().Context(), payload.Identity.Id); err != nil {
		h.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		h.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if _, err := h.store.User().Create(c.Request().Context(), dto.UserChangeset{ID: opt.New(identityID)}); err != nil {
		h.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	gameTimeCategories, err := h.store.GameTimeCategory().List(c.Request().Context())
	if err != nil {
		h.Log.Error("failed to create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ratingsToInsert := make([]dto.RatingChangeset, 0, len(gameTimeCategories))
	for _, tc := range gameTimeCategories {
		ratingsToInsert = append(ratingsToInsert, dto.RatingChangeset{
			UserID:             opt.New(identityID),
			GameTimeCategoryID: opt.New(tc.ID),
			Glicko:             opt.New[int32](1500),
			Glicko2:            opt.New[int32](1500),
		})
	}
	if _, err := h.store.Rating().BatchCreate(c.Request().Context(), ratingsToInsert); err != nil {
		h.Log.Error("failed to batch create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// send welcome email...

	return c.JSON(http.StatusOK, nil)
}

func (h *ApiHandler) registrationAfterOidc(c echo.Context) error {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := c.Bind(&payload); err != nil {
		h.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_oidc"), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	h.Log.Debug("kratos webhook registration_after_oidc", slog.String("identity_id", payload.Identity.Id))
	v := c.Request().Header.Get("Authorization")
	if v != h.Cfg.KratosAPIKey {
		h.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_oidc"), slog.String("authorization", v))
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err := h.createUserRelationTuples(c.Request().Context(), payload.Identity.Id); err != nil {
		h.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		h.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if _, err := h.store.User().Create(c.Request().Context(), dto.UserChangeset{ID: opt.New(identityID)}); err != nil {
		h.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}
	gameTimeCategories, err := h.store.GameTimeCategory().List(c.Request().Context())
	if err != nil {
		h.Log.Error("failed to create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ratingsToInsert := make([]dto.RatingChangeset, 0, len(gameTimeCategories))
	for _, tc := range gameTimeCategories {
		ratingsToInsert = append(ratingsToInsert, dto.RatingChangeset{
			UserID:             opt.New(identityID),
			GameTimeCategoryID: opt.New(tc.ID),
			Glicko:             opt.New[int32](1500),
			Glicko2:            opt.New[int32](1500),
		})
	}
	if _, err := h.store.Rating().BatchCreate(c.Request().Context(), ratingsToInsert); err != nil {
		h.Log.Error("failed to batch create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// send welcome email...

	return c.JSON(http.StatusOK, nil)
}

func (h *ApiHandler) createUserRelationTuples(ctx context.Context, identityID string) error {
	_, err := h.Keto.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customers",
					Relation:  "members",
					Subject:   rts.NewSubjectID(identityID),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    identityID,
					Relation:  "owners",
					Subject:   rts.NewSubjectID(identityID),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    identityID,
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities/", ""),
				},
			},
		},
	})
	return err
}
