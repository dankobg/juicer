package server

import (
	"log/slog"
	"net/http"

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

	if _, err := h.Keto.Write.TransactRelationTuples(c.Request().Context(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customers",
					Relation:  "members",
					Subject:   rts.NewSubjectID(payload.Identity.Id),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    payload.Identity.Id,
					Relation:  "owners",
					Subject:   rts.NewSubjectID(payload.Identity.Id),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    payload.Identity.Id,
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities/", ""),
				},
			},
		},
	}); err != nil {
		h.Log.Error("failed to insert user relation-tuple", slog.String("user_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
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

	if _, err := h.Keto.Write.TransactRelationTuples(c.Request().Context(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customers",
					Relation:  "members",
					Subject:   rts.NewSubjectID(payload.Identity.Id),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    payload.Identity.Id,
					Relation:  "owners",
					Subject:   rts.NewSubjectID(payload.Identity.Id),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    payload.Identity.Id,
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities/", ""),
				},
			},
		},
	}); err != nil {
		h.Log.Error("failed to insert user relation-tuple", slog.String("user_id", payload.Identity.Id), slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, nil)
	}

	// send welcome email...

	return c.JSON(http.StatusOK, nil)
}
