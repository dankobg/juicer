package webhooks

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/features/idp"
	orykratos "github.com/ory/client-go"
)

type Webhooks struct {
	idp       *idp.IdentityProvider
	userPst   idp.UserPersistor
	gtcPst    game.GameTimeCategoryPersistor
	ratingPst game.RatingPersistor
	log       *slog.Logger
}

func NewWebhooks(idp *idp.IdentityProvider, l *slog.Logger) *Webhooks {
	return &Webhooks{
		idp: idp,
		log: l,
	}
}

func (wh *Webhooks) RegistrationAfterPassword(w http.ResponseWriter, r *http.Request) {
	wh.handleRegistrationWebhook(w, r, "registration_after_password")
}

func (wh *Webhooks) RegistrationAfterOidc(w http.ResponseWriter, r *http.Request) {
	wh.handleRegistrationWebhook(w, r, "registration_after_oidc")
}

// Shared helper to keep your code DRY
func (wh *Webhooks) handleRegistrationWebhook(w http.ResponseWriter, r *http.Request, webhookName string) {
	if !wh.idp.ValidateKratosWebhookSecret(r.Header.Get("Authorization")) {
		wh.log.Error("webhook auth failed", slog.String("webhook", webhookName))
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return
	}

	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		wh.log.Error("invalid webhook payload", slog.String("webhook", webhookName), slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	wh.log.Debug("kratos webhook received", slog.String("webhook", webhookName), slog.String("identity_id", payload.Identity.Id))

	if err := wh.idp.OnUserRegistered(r.Context(), payload.Identity.Id); err != nil {
		wh.log.Error("failed to process user registration webhook", slog.String("webhook", webhookName), slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		// @TODO: better status code
		http.Error(w, "internal processing error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
