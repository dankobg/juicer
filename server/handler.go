package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/auth/keto"
	"github.com/dankobg/juicer/auth/kratos"
	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/mailer"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/ws"
	"github.com/redis/go-redis/v9"
)

// var _ api.StrictServerInterface = (*ApiHandler)(nil)

type categoryThreshold struct {
	upperLimit   time.Duration
	timeCategory pb.GameTimeCategory
}

// for now this does not change, so i keep it static and fetch once
type protoMappingsCache struct {
	gameVariantsProtoToDB       map[pb.GameVariant]int64
	gameTimeKindsProtoToDB      map[pb.GameTimeKind]int64
	gameTimeCategoriesProtoToDB map[pb.GameTimeCategory]int64
	gameResultsProtoToDB        map[pb.GameResult]int64
	gameResultStatusesProtoToDB map[pb.GameResultStatus]int64
	gameStatesProtoToDB         map[pb.GameState]int64

	gameVariantsDBToProto       map[int64]pb.GameVariant
	gameTimeKindsDBToProto      map[int64]pb.GameTimeKind
	gameTimeCategoriesDBToProto map[int64]pb.GameTimeCategory
	gameResultsDBToProto        map[int64]pb.GameResult
	gameResultStatusesDBToProto map[int64]pb.GameResultStatus
	gameStatesDBToProto         map[int64]pb.GameState
}

func newProtoMappingsCache() protoMappingsCache {
	return protoMappingsCache{
		gameVariantsProtoToDB:       make(map[pb.GameVariant]int64),
		gameTimeKindsProtoToDB:      make(map[pb.GameTimeKind]int64),
		gameTimeCategoriesProtoToDB: make(map[pb.GameTimeCategory]int64),
		gameResultsProtoToDB:        make(map[pb.GameResult]int64),
		gameResultStatusesProtoToDB: make(map[pb.GameResultStatus]int64),
		gameStatesProtoToDB:         make(map[pb.GameState]int64),

		gameVariantsDBToProto:       make(map[int64]pb.GameVariant),
		gameTimeKindsDBToProto:      make(map[int64]pb.GameTimeKind),
		gameTimeCategoriesDBToProto: make(map[int64]pb.GameTimeCategory),
		gameResultsDBToProto:        make(map[int64]pb.GameResult),
		gameResultStatusesDBToProto: make(map[int64]pb.GameResultStatus),
		gameStatesDBToProto:         make(map[int64]pb.GameState),
	}
}

type ApiHandler struct {
	Cfg        *config.Config
	Log        *slog.Logger
	Kratos     *kratos.Client
	Keto       *keto.Client
	Hub        *ws.Hub
	Rdb        *redis.Client
	persistor  persistence.Persistor
	Mailer     mailer.Mailer
	openapiTpl *template.Template
	bus        *bus

	categoryThresholds []categoryThreshold
	protoMappingsCache protoMappingsCache
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	rdb *redis.Client,
	kratos *kratos.Client,
	keto *keto.Client,
	mailer mailer.Mailer,
	hub *ws.Hub,
	p persistence.Persistor,
) *ApiHandler {
	apiHandler := &ApiHandler{
		Cfg:                cfg,
		Log:                log,
		Kratos:             kratos,
		Keto:               keto,
		persistor:          p,
		Mailer:             mailer,
		Hub:                hub,
		Rdb:                rdb,
		bus:                newBus(rdb),
		categoryThresholds: make([]categoryThreshold, 0),
		protoMappingsCache: newProtoMappingsCache(),
	}

	return apiHandler
}

func (a *ApiHandler) SetOpenapiTemplates(tpl *template.Template) {
	a.openapiTpl = tpl
}

func newNotFoundErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
	if len(reason) > 0 {
		e.Reason = new(reason[0])
	}

	return e
}

func newUnauthenticatedErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = new(reason[0])
	}

	return e
}

func newUnauthorizedErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = new(reason[0])
	}

	return e
}

func newGenericErr(statusCode int32, code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: statusCode,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = new(reason[0])
	}

	return e
}

func newNotFoundResp(code, message string, reason ...string) api.NotFoundErrorResponseJSONResponse {
	e := newNotFoundErr(code, message, reason...)
	return api.NotFoundErrorResponseJSONResponse(e)
}

func newUnauthenticatedResp(code, message string, reason ...string) api.UnauthenticatedErrorResponse {
	e := newUnauthenticatedErr(code, message, reason...)
	return api.UnauthenticatedErrorResponse(e)
}

func newUnauthorizedResp(code, message string, reason ...string) api.UnauthorizedErrorResponseJSONResponse {
	e := newUnauthorizedErr(code, message, reason...)
	return api.UnauthorizedErrorResponseJSONResponse(e)
}

func newGenericResp(statusCode int32, code, message string, reason ...string) api.GenericErrorResponseJSONResponse {
	e := newGenericErr(statusCode, code, message, reason...)
	return api.GenericErrorResponseJSONResponse(e)
}
