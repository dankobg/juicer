package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListSessions(ctx context.Context, request api.ListSessionsRequestObject) (api.ListSessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListSessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "permission denied")}, nil
	}

	pagedSessions, err := a.idp.ListSessions(ctx, request)
	if err != nil {
		return api.ListSessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "sessions_list", "failed to list sessions")}, nil
	}

	resp := api.ListSessions200JSONResponse(pagedSessions.Data)

	return resp, nil
}

func (a *ApiHandler) GetSession(ctx context.Context, request api.GetSessionRequestObject) (api.GetSessionResponseObject, error) {
	sess := GetSession(ctx)

	session, err := a.idp.GetSession(ctx, request)
	if err != nil {
		return api.GetSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_get", "failed to get session")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    shared.AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "permission denied")}, nil
	}

	return api.GetSession200JSONResponse(*session), nil
}

func (a *ApiHandler) DisableSession(ctx context.Context, request api.DisableSessionRequestObject) (api.DisableSessionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    shared.AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DisableSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "permission denied")}, nil
	}

	if err := a.idp.DisableSession(ctx, request); err != nil {
		return api.DisableSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_disable", "failed to disable session")}, nil
	}

	return api.DisableSession204Response{}, nil
}

func (a *ApiHandler) ExtendSession(ctx context.Context, request api.ExtendSessionRequestObject) (api.ExtendSessionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    shared.AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ExtendSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "permission denied")}, nil
	}

	session, err := a.idp.ExtendSession(ctx, request)
	if err != nil {
		return api.ExtendSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_extend", "failed to extend session")}, nil
	}

	return api.ExtendSession200JSONResponse(*session), nil
}
