package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListIdentitySchemas(ctx context.Context, request api.ListIdentitySchemasRequestObject) (api.ListIdentitySchemasResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Schemas",
			Object:    "schemas",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListIdentitySchemas403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("schema_permission", "permission denied")}, nil
	}

	schemaContainers, err := a.idp.ListIdentitySchemas(ctx, request)
	if err != nil {
		return api.ListIdentitySchemasdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "schemas_list", "failed to list identity schemas")}, nil
	}

	out := api.ListIdentitySchemas200JSONResponse(schemaContainers.Data)

	return out, nil
}

func (a *ApiHandler) GetIdentitySchema(ctx context.Context, request api.GetIdentitySchemaRequestObject) (api.GetIdentitySchemaResponseObject, error) {
	sess := GetSession(ctx)

	identitySchema, err := a.idp.GetIdentitySchema(ctx, request)
	if err != nil {
		return api.GetIdentitySchema404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("schema_not_found", "schema not found")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Schema",
			Object:    shared.AuthzSchemaID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetIdentitySchema403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("schema_permission", "permission denied")}, nil
	}

	return api.GetIdentitySchema200JSONResponse(*identitySchema), nil
}
