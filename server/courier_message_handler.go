package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListCourierMessages(ctx context.Context, request api.ListCourierMessagesRequestObject) (api.ListCourierMessagesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "CourierMessages",
			Object:    "courier_messages",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListCourierMessages403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("message_permission", "permission denied")}, nil
	}

	courierMessages, err := a.idp.ListCourierMessages(ctx, request)
	if err != nil {
		return api.ListCourierMessagesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "message_list", "failed to list courier messages")}, nil
	}

	out := api.ListCourierMessages200JSONResponse(courierMessages.Data)

	return out, nil
}

func (a *ApiHandler) GetCourierMessage(ctx context.Context, request api.GetCourierMessageRequestObject) (api.GetCourierMessageResponseObject, error) {
	sess := GetSession(ctx)

	courierMessage, err := a.idp.GetCourierMessage(ctx, request)
	if err != nil {
		return api.GetCourierMessage404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("message_not_found", "courier message not found")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "CourierMessage",
			// Object:    shared.AuthzCourierMessageID(request.ID),
			Namespace: "CourierMessages",
			Object:    "courier_messages",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetCourierMessage403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("message_permission", "permission denied")}, nil
	}

	out := api.GetCourierMessage200JSONResponse(*courierMessage)

	return out, nil
}
