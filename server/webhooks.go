package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
	orykratos "github.com/ory/client-go"
)

func (a *ApiHandler) registrationAfterPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_password"), slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	a.Log.Debug("kratos webhook registration_after_password", slog.String("identity_id", payload.Identity.Id))

	v := r.Header.Get("Authorization")
	if v != a.Cfg.KratosAPIKey {
		a.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_password"), slog.String("authorization", v))
		http.Error(w, "unauthorized", http.StatusBadRequest)

		return
	}

	if err := createUserRelationTuples(r.Context(), a.Keto, payload.Identity.Id); err != nil {
		a.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to insert user relation-tuple", http.StatusBadRequest)

		return
	}

	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		a.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if _, err := a.persistor.User().CreateUser(r.Context(), models.UserSetter{ID: omit.From(identityID)}); err != nil {
		a.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to create user", http.StatusBadRequest)

		return
	}

	gameTimeCategories, err := a.persistor.GameTimeCategory().ListGameTimeCategories(r.Context(), dbtype.ListGameTimeCategoriesFilters{})
	if err != nil {
		a.Log.Error("failed to create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	ratingsToInsert := make([]models.RatingSetter, 0, len(gameTimeCategories.Data))
	for _, tc := range gameTimeCategories.Data {
		ratingsToInsert = append(ratingsToInsert, models.RatingSetter{
			UserID:             omit.From(identityID),
			GameTimeCategoryID: omit.From(tc.ID),
			Glicko:             omit.From[int32](1500),
			Glicko2:            omit.From[int32](1500),
		})
	}

	if _, err := a.persistor.Rating().BulkCreateRatings(r.Context(), ratingsToInsert); err != nil {
		a.Log.Error("failed to bulk create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	// send welcome email...

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{}`))
}

func (a *ApiHandler) registrationAfterOidc(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_oidc"), slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	a.Log.Debug("kratos webhook registration_after_oidc", slog.String("identity_id", payload.Identity.Id))

	v := r.Header.Get("Authorization")
	if v != a.Cfg.KratosAPIKey {
		a.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_oidc"), slog.String("authorization", v))
		http.Error(w, "unauthorized", http.StatusBadRequest)

		return
	}

	if err := createUserRelationTuples(r.Context(), a.Keto, payload.Identity.Id); err != nil {
		a.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to insert user relation-tuple", http.StatusBadRequest)

		return
	}

	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		a.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if _, err := a.persistor.User().CreateUser(r.Context(), models.UserSetter{ID: omit.From(identityID)}); err != nil {
		a.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to create user", http.StatusBadRequest)

		return
	}

	gameTimeCategories, err := a.persistor.GameTimeCategory().ListGameTimeCategories(r.Context(), dbtype.ListGameTimeCategoriesFilters{})
	if err != nil {
		a.Log.Error("failed to create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	ratingsToInsert := make([]models.RatingSetter, 0, len(gameTimeCategories.Data))
	for _, tc := range gameTimeCategories.Data {
		ratingsToInsert = append(ratingsToInsert, models.RatingSetter{
			UserID:             omit.From(identityID),
			GameTimeCategoryID: omit.From(tc.ID),
			Glicko:             omit.From[int32](1500),
			Glicko2:            omit.From[int32](1500),
		})
	}

	if _, err := a.persistor.Rating().BulkCreateRatings(r.Context(), ratingsToInsert); err != nil {
		a.Log.Error("failed to bulk create ratings", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	// send welcome email...

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{}`))
}
