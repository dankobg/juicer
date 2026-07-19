package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListIdentities(ctx context.Context, request api.ListIdentitiesRequestObject) (api.ListIdentitiesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	identities, err := a.idp.ListIdentities(ctx, request)
	if err != nil {
		return api.ListIdentitiesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identities_list", "failed to list identities")}, nil
	}

	out := api.ListIdentities200JSONResponse(identities.Data)

	return out, nil
}

func (a *ApiHandler) GetIdentity(ctx context.Context, request api.GetIdentityRequestObject) (api.GetIdentityResponseObject, error) {
	sess := GetSession(ctx)

	identity, err := a.idp.GetIdentity(ctx, request)
	if err != nil {
		return api.GetIdentity404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("identity_not_found", "identity not found")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	return api.GetIdentity200JSONResponse(*identity), nil
}

func (a *ApiHandler) CreateIdentity(ctx context.Context, request api.CreateIdentityRequestObject) (api.CreateIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	identity, err := a.idp.CreateIdentity(ctx, request)
	if err != nil {
		return nil, err
	}

	return api.CreateIdentity201JSONResponse(*identity), nil
}

func (a *ApiHandler) UpdateIdentity(ctx context.Context, request api.UpdateIdentityRequestObject) (api.UpdateIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	identity, err := a.idp.UpdateIdentity(ctx, request)
	if err != nil {
		return api.UpdateIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_edit", "failed to edit identity")}, nil
	}

	return api.UpdateIdentity200JSONResponse(*identity), nil
}

func (a *ApiHandler) DeleteIdentity(ctx context.Context, request api.DeleteIdentityRequestObject) (api.DeleteIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	if err := a.idp.DeleteIdentity(ctx, request); err != nil {
		return api.DeleteIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete", "failed to delete identity")}, nil
	}

	return api.DeleteIdentity204Response{}, nil
}

func (a *ApiHandler) PatchIdentity(ctx context.Context, request api.PatchIdentityRequestObject) (api.PatchIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.PatchIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	identity, err := a.idp.PatchIdentity(ctx, request)
	if err != nil {
		return api.PatchIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_patch", "failed to patch identity")}, nil
	}

	return api.PatchIdentity200JSONResponse(*identity), nil
}

func (a *ApiHandler) BatchPatchIdentities(ctx context.Context, request api.BatchPatchIdentitiesRequestObject) (api.BatchPatchIdentitiesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.BatchPatchIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	batchPatchIdentities, err := a.idp.BatchPatchIdentities(ctx, request)
	if err != nil {
		return api.BatchPatchIdentitiesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_batch_patch", "failed to batch patch identity")}, nil
	}

	return api.BatchPatchIdentities200JSONResponse(*batchPatchIdentities), nil
}

func (a *ApiHandler) DeleteIdentityCredentials(ctx context.Context, request api.DeleteIdentityCredentialsRequestObject) (api.DeleteIdentityCredentialsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentityCredentials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	if err := a.idp.DeleteIdentityCredentials(ctx, request); err != nil {
		return api.DeleteIdentityCredentialsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_credentials", "failed to delete identity")}, nil
	}

	return api.DeleteIdentityCredentials204Response{}, nil
}

func (a *ApiHandler) DeleteIdentitySessions(ctx context.Context, request api.DeleteIdentitySessionsRequestObject) (api.DeleteIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	if err := a.idp.DeleteIdentitySessions(ctx, request); err != nil {
		return api.DeleteIdentitySessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_sessions", "failed to delete identity sessions")}, nil
	}

	return api.DeleteIdentitySessions204Response{}, nil
}

func (a *ApiHandler) ListIdentitySessions(ctx context.Context, request api.ListIdentitySessionsRequestObject) (api.ListIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	identitySessions, err := a.idp.ListIdentitySessions(ctx, request)
	if err != nil {
		return api.ListIdentitySessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identities_list", "failed to list identities")}, nil
	}

	return api.ListIdentitySessions200JSONResponse(identitySessions.Data), nil
}

func (a *ApiHandler) CreateRecoveryCodeForIdentity(ctx context.Context, request api.CreateRecoveryCodeForIdentityRequestObject) (api.CreateRecoveryCodeForIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateRecoveryCodeForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	recoveryCodeForIdentity, err := a.idp.CreateRecoveryCodeForIdentity(ctx, request)
	if err != nil {
		return api.CreateRecoveryCodeForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_code", "failed to create recovery code for identity")}, nil
	}

	return api.CreateRecoveryCodeForIdentity201JSONResponse(*recoveryCodeForIdentity), nil
}

func (a *ApiHandler) CreateRecoveryLinkForIdentity(ctx context.Context, request api.CreateRecoveryLinkForIdentityRequestObject) (api.CreateRecoveryLinkForIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateRecoveryLinkForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	recoveryLinkForIdentity, err := a.idp.CreateRecoveryLinkForIdentity(ctx, request)
	if err != nil {
		return api.CreateRecoveryLinkForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_link", "failed to create recovery link for identity")}, nil
	}

	return api.CreateRecoveryLinkForIdentity200JSONResponse(*recoveryLinkForIdentity), nil
}

func (a *ApiHandler) CreateFriendRequest(ctx context.Context, request api.CreateFriendRequestRequestObject) (api.CreateFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API CreateFriendRequest")
}

func (a *ApiHandler) ListFriendRequests(ctx context.Context, request api.ListFriendRequestsRequestObject) (api.ListFriendRequestsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFriendRequests")
}

func (a *ApiHandler) AcceptFriendRequest(ctx context.Context, request api.AcceptFriendRequestRequestObject) (api.AcceptFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API AcceptFriendRequest")
}

func (a *ApiHandler) DeclineFriendRequest(ctx context.Context, request api.DeclineFriendRequestRequestObject) (api.DeclineFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeclineFriendRequest")
}

func (a *ApiHandler) CancelFriendRequest(ctx context.Context, request api.CancelFriendRequestRequestObject) (api.CancelFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API CancelFriendRequest")
}

func (a *ApiHandler) ListFriends(ctx context.Context, request api.ListFriendsRequestObject) (api.ListFriendsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFriends")
}

func (a *ApiHandler) GetFriend(ctx context.Context, request api.GetFriendRequestObject) (api.GetFriendResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetFriend")
}

func (a *ApiHandler) DeleteFriend(ctx context.Context, request api.DeleteFriendRequestObject) (api.DeleteFriendResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeleteFriend")
}

func (a *ApiHandler) DeleteFriends(ctx context.Context, request api.DeleteFriendsRequestObject) (api.DeleteFriendsResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeleteFriends")
}

func (a *ApiHandler) ListFollowings(ctx context.Context, request api.ListFollowingsRequestObject) (api.ListFollowingsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFollowings")
}

func (a *ApiHandler) GetFollowing(ctx context.Context, request api.GetFollowingRequestObject) (api.GetFollowingResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetFollowing")
}

func (a *ApiHandler) FollowUser(ctx context.Context, request api.FollowUserRequestObject) (api.FollowUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API FollowUser")
}

func (a *ApiHandler) UnfollowUser(ctx context.Context, request api.UnfollowUserRequestObject) (api.UnfollowUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnfollowUser")
}

func (a *ApiHandler) UnfollowUsers(ctx context.Context, request api.UnfollowUsersRequestObject) (api.UnfollowUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnfollowUsers")
}

func (a *ApiHandler) ListBlockedUsers(ctx context.Context, request api.ListBlockedUsersRequestObject) (api.ListBlockedUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListBlockedUsers")
}

func (a *ApiHandler) GetBlockedUser(ctx context.Context, request api.GetBlockedUserRequestObject) (api.GetBlockedUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetBlockedUser")
}

func (a *ApiHandler) BlockUser(ctx context.Context, request api.BlockUserRequestObject) (api.BlockUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API BlockUser")
}

func (a *ApiHandler) UnblockUser(ctx context.Context, request api.UnblockUserRequestObject) (api.UnblockUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnblockUser")
}

func (a *ApiHandler) UnblockUsers(ctx context.Context, request api.UnblockUsersRequestObject) (api.UnblockUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnblockUsers")
}
