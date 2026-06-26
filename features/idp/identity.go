package idp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/auth/keto"
	"github.com/dankobg/juicer/auth/kratos"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/shared"
	"github.com/google/uuid"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Paged[T any] struct {
	Data []T
}

type IdentityProvider struct {
	kratos       *kratos.Client
	keto         *keto.Client
	kratosAPIKey string
	ketoAPIKey   string
	userPst      UserPersistor
	// meh (i know, does not belong here) but who cares
	gtcPst    game.GameTimeCategoryPersistor
	ratingPst game.RatingPersistor
	log       *slog.Logger
}

func NewIdentityProvider(
	kratos *kratos.Client,
	keto *keto.Client,
	kratosAPIKey, ketoAPIKey string,
	userPst UserPersistor,
	gtcPst game.GameTimeCategoryPersistor,
	ratingPst game.RatingPersistor,
	l *slog.Logger,
) *IdentityProvider {
	return &IdentityProvider{
		kratos:       kratos,
		keto:         keto,
		kratosAPIKey: kratosAPIKey,
		ketoAPIKey:   ketoAPIKey,
		userPst:      userPst,
		gtcPst:       gtcPst,
		ratingPst:    ratingPst,
		log:          l,
	}
}

func (idp *IdentityProvider) ValidateKratosWebhookSecret(authorization string) bool {
	return authorization == idp.kratosAPIKey
}

func (idp *IdentityProvider) OnUserRegistered(ctx context.Context, identityID string) error {
	if err := createUserRelationTuples(ctx, idp.keto, identityID); err != nil {
		return fmt.Errorf("failed to insert user relation-tuple: %w", err)
	}

	uid, err := uuid.Parse(identityID)
	if err != nil {
		return fmt.Errorf("failed to parse identity id")
	}

	if _, err := idp.userPst.CreateUser(ctx, models.UserSetter{ID: omit.From(uid)}); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	gameTimeCategories, err := idp.gtcPst.ListGameTimeCategories(ctx, game.ListGameTimeCategoriesFilters{})
	if err != nil {
		return fmt.Errorf("failed to list game time categories: %w", err)
	}

	ratingsToInsert := make([]models.RatingSetter, 0, len(gameTimeCategories.Data))
	for _, tc := range gameTimeCategories.Data {
		ratingsToInsert = append(ratingsToInsert, models.RatingSetter{
			UserID:             omit.From(uid),
			GameTimeCategoryID: omit.From(tc.ID),
			Glicko:             omit.From[int32](1500),
			Glicko2:            omit.From[int32](1500),
		})
	}

	if _, err := idp.ratingPst.BulkCreateRatings(ctx, ratingsToInsert); err != nil {
		return fmt.Errorf("failed to bulk create initial user ratings: %w", err)
	}

	// send email...

	return nil
}

func (idp *IdentityProvider) ListIdentitySchemas(ctx context.Context, request api.ListIdentitySchemasRequestObject) (Paged[api.IdentitySchemaContainer], error) {
	req := idp.kratos.Admin.IdentityAPI.ListIdentitySchemas(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	schemaContainers, schemaContainersResp, err := req.Execute()
	if err != nil {
		return Paged[api.IdentitySchemaContainer]{}, nil
	}

	defer func() { _ = schemaContainersResp.Body.Close() }()

	outSchemas := make([]api.IdentitySchemaContainer, 0, len(schemaContainers))
	for _, sc := range schemaContainers {
		res, err := SchemaContainerToResponse(sc)
		if err != nil {
			return Paged[api.IdentitySchemaContainer]{}, err
		}

		outSchemas = append(outSchemas, res)
	}

	out := Paged[api.IdentitySchemaContainer]{
		Data: outSchemas,
	}

	return out, nil
}

func (idp *IdentityProvider) GetIdentitySchema(ctx context.Context, request api.GetIdentitySchemaRequestObject) (*api.IdentitySchema, error) {
	req := idp.kratos.Admin.IdentityAPI.GetIdentitySchema(ctx, request.ID)

	identitySchema, identitySchemaResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identitySchemaResp.Body.Close() }()

	return &identitySchema, nil
}

func (idp *IdentityProvider) ListIdentities(ctx context.Context, request api.ListIdentitiesRequestObject) (Paged[api.Identity], error) {
	req := idp.kratos.Admin.IdentityAPI.ListIdentities(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	if request.Params.Ids != nil {
		req = req.Ids(*request.Params.Ids)
	}

	if request.Params.Consistency != nil {
		req = req.Consistency(string(*request.Params.Consistency))
	}

	if request.Params.CredentialsIdentifier != nil {
		req = req.CredentialsIdentifier(*request.Params.CredentialsIdentifier)
	}

	if request.Params.IncludeCredential != nil {
		req = req.IncludeCredential(*request.Params.IncludeCredential)
	}

	if request.Params.PreviewCredentialsIdentifierSimilar != nil {
		req = req.PreviewCredentialsIdentifierSimilar(*request.Params.PreviewCredentialsIdentifierSimilar)
	}

	identities, identitiesResp, err := req.Execute()
	if err != nil {
		return Paged[api.Identity]{}, err
	}

	defer func() { _ = identitiesResp.Body.Close() }()

	outIdentities := make([]api.Identity, 0, len(identities))
	for _, identity := range identities {
		res, err := IdentityToResponse(identity)
		if err != nil {
			return Paged[api.Identity]{}, err
		}

		outIdentities = append(outIdentities, res)
	}

	out := Paged[api.Identity]{
		Data: outIdentities,
	}

	return out, nil
}

func (idp *IdentityProvider) GetIdentity(ctx context.Context, request api.GetIdentityRequestObject) (*api.Identity, error) {
	req := idp.kratos.Admin.IdentityAPI.GetIdentity(ctx, request.ID)
	if request.Params.IncludeCredential != nil {
		includeParams := make([]string, 0, len(*request.Params.IncludeCredential))
		for _, iparam := range *request.Params.IncludeCredential {
			includeParams = append(includeParams, string(iparam))
		}

		req = req.IncludeCredential(includeParams)
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (idp *IdentityProvider) CreateIdentity(ctx context.Context, request api.CreateIdentityRequestObject) (*api.Identity, error) {
	req := idp.kratos.Admin.IdentityAPI.CreateIdentity(ctx)

	if request.Body != nil {
		var credentials *orykratos.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = orykratos.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = orykratos.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = orykratos.NewIdentityWithCredentialsPasswordConfig()
				}

				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}

				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}

			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = orykratos.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = orykratos.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]orykratos.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, orykratos.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}

						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}

		recoveryAddresses := make([]orykratos.RecoveryIdentityAddress, 0)

		if request.Body.RecoveryAddresses != nil {
			for _, recAddr := range *request.Body.RecoveryAddresses {
				recoveryAddresses = append(recoveryAddresses, orykratos.RecoveryIdentityAddress{
					Id:        new(recAddr.ID.String()),
					Value:     recAddr.Value,
					Via:       recAddr.Via,
					CreatedAt: recAddr.CreatedAt,
					UpdatedAt: recAddr.UpdatedAt,
				})
			}
		}

		verifiableAddresses := make([]orykratos.VerifiableIdentityAddress, 0)

		if request.Body.VerifiableAddresses != nil {
			for _, verAddr := range *request.Body.VerifiableAddresses {
				var id *string
				if verAddr.ID != nil {
					id = new(verAddr.ID.String())
				}

				verifiableAddress := orykratos.VerifiableIdentityAddress{
					Id:        id,
					Status:    verAddr.Status,
					Value:     verAddr.Value,
					Verified:  verAddr.Verified,
					Via:       string(verAddr.Via),
					CreatedAt: verAddr.CreatedAt,
					UpdatedAt: verAddr.UpdatedAt,
				}
				if verAddr.VerifiedAt.IsSpecified() && !verAddr.VerifiedAt.IsNull() {
					verifiableAddress.VerifiedAt = new(verAddr.VerifiedAt.MustGet().UTC())
				}

				verifiableAddresses = append(verifiableAddresses, verifiableAddress)
			}
		}

		req = req.CreateIdentityBody(orykratos.CreateIdentityBody{
			Credentials:         credentials,
			MetadataAdmin:       request.Body.MetadataAdmin,
			MetadataPublic:      request.Body.MetadataPublic,
			RecoveryAddresses:   recoveryAddresses,
			SchemaId:            request.Body.SchemaID,
			State:               (*string)(request.Body.State),
			Traits:              request.Body.Traits,
			VerifiableAddresses: verifiableAddresses,
		})
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	identityID, err := uuid.Parse(identity.Id)
	if err != nil {
		return nil, err
	}

	if _, err := idp.userPst.CreateUser(ctx, models.UserSetter{ID: omit.From(identityID)}); err != nil {
		return nil, err
	}

	// @TODO: use outbox pattern
	if err := createUserRelationTuples(ctx, idp.keto, identity.Id); err != nil {
		idp.log.Error("failed to insert identity relation-tuple", slog.String("identity_id", identityID.String()), slog.Any("error", err))
		return nil, err
	}

	return &resp, nil
}

func (idp *IdentityProvider) UpdateIdentity(ctx context.Context, request api.UpdateIdentityRequestObject) (*api.Identity, error) {
	req := idp.kratos.Admin.IdentityAPI.UpdateIdentity(ctx, request.ID)
	if request.Body != nil {
		var credentials *orykratos.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = orykratos.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = orykratos.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = orykratos.NewIdentityWithCredentialsPasswordConfig()
				}

				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}

				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}

			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = orykratos.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = orykratos.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]orykratos.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, orykratos.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}

						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}

		req = req.UpdateIdentityBody(orykratos.UpdateIdentityBody{
			Credentials:    credentials,
			MetadataAdmin:  request.Body.MetadataAdmin,
			MetadataPublic: request.Body.MetadataPublic,
			SchemaId:       request.Body.SchemaID,
			State:          string(request.Body.State),
			Traits:         request.Body.Traits,
		})
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (idp *IdentityProvider) DeleteIdentity(ctx context.Context, request api.DeleteIdentityRequestObject) error {
	req := idp.kratos.Admin.IdentityAPI.DeleteIdentity(ctx, request.ID)

	deleteIdentityResp, err := req.Execute()
	if err != nil {
		return err
	}

	defer func() { _ = deleteIdentityResp.Body.Close() }()

	// @TODO: use outbox pattern
	if err := deleteUserRelationTuples(ctx, idp.keto, request.ID); err != nil {
		idp.log.Error("failed to delete identity relation-tuple", slog.String("identity_id", request.ID), slog.Any("error", err))
		return err
	}

	return nil
}

func (idp *IdentityProvider) PatchIdentity(ctx context.Context, request api.PatchIdentityRequestObject) (*api.Identity, error) {
	req := idp.kratos.Admin.IdentityAPI.PatchIdentity(ctx, request.ID)
	if request.Body != nil {
		patches := make([]orykratos.JsonPatch, 0, len(*request.Body))
		for _, x := range *request.Body {
			patches = append(patches, orykratos.JsonPatch{
				From:  x.From,
				Op:    string(x.Op),
				Path:  x.Path,
				Value: x.Value,
			})
		}

		req = req.JsonPatch(patches)
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	// @TODO: delete tuples if patch OP == "remove"

	return &resp, nil
}

func (idp *IdentityProvider) BatchPatchIdentities(ctx context.Context, request api.BatchPatchIdentitiesRequestObject) (*api.BatchPatchIdentitiesResponse, error) {
	req := idp.kratos.Admin.IdentityAPI.BatchPatchIdentities(ctx)

	if request.Body != nil {
		patch := orykratos.PatchIdentitiesBody{}
		req = req.PatchIdentitiesBody(patch)
	}

	batchPatchIdentities, batchPatchIdentitiesResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = batchPatchIdentitiesResp.Body.Close() }()

	identitiesPatches := make([]api.IdentityPatchResponse, 0)

	for _, x := range batchPatchIdentities.Identities {
		var (
			defaultErr   any
			identityUUID *api.UUID
		)

		if x.Identity != nil {
			identityUUIDParsed, err := uuid.Parse(*x.Identity)
			if err != nil {
				defaultErr = err
			} else {
				identityUUID = new(identityUUIDParsed)
			}
		}

		var patchUUID *api.UUID

		if x.PatchId != nil {
			patchUUIDParsed, err := uuid.Parse(*x.PatchId)
			if err != nil {
				defaultErr = err
			} else {
				patchUUID = new(patchUUIDParsed)
			}
		}

		identitiesPatches = append(identitiesPatches, api.IdentityPatchResponse{
			Action:   (*api.IdentityPatchResponseAction)(x.Action),
			Identity: identityUUID,
			PatchID:  patchUUID,
			Error:    &defaultErr,
		})
	}

	resp := api.BatchPatchIdentitiesResponse{
		Identities: &identitiesPatches,
	}

	// @TODO: delete tuples if patch OP == "remove"

	return &resp, nil
}

func (idp *IdentityProvider) DeleteIdentityCredentials(ctx context.Context, request api.DeleteIdentityCredentialsRequestObject) error {
	req := idp.kratos.Admin.IdentityAPI.DeleteIdentityCredentials(ctx, request.ID, string(request.Type))
	if request.Params.Identifier != nil {
		req = req.Identifier(*request.Params.Identifier)
	}

	deleteIdentityCredentialsResp, err := req.Execute()
	if err != nil {
		return err
	}

	defer func() { _ = deleteIdentityCredentialsResp.Body.Close() }()

	return nil
}

func (idp *IdentityProvider) DeleteIdentitySessions(ctx context.Context, request api.DeleteIdentitySessionsRequestObject) error {
	req := idp.kratos.Admin.IdentityAPI.DeleteIdentitySessions(ctx, request.ID)

	deleteIdentitySessionsResp, err := req.Execute()
	if err != nil {
		return err
	}

	defer func() { _ = deleteIdentitySessionsResp.Body.Close() }()

	return nil
}

func (idp *IdentityProvider) ListIdentitySessions(ctx context.Context, request api.ListIdentitySessionsRequestObject) (Paged[api.Session], error) {
	req := idp.kratos.Admin.IdentityAPI.ListIdentitySessions(ctx, request.ID)
	if request.Params.Active != nil {
		req = req.Active(*request.Params.Active)
	}

	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	identitySessions, identitySessionsResp, err := req.Execute()
	if err != nil {
		return Paged[api.Session]{}, err
	}

	defer func() { _ = identitySessionsResp.Body.Close() }()

	outSessions := make([]api.Session, 0, len(identitySessions))
	for _, sess := range identitySessions {
		res, err := SessionToResponse(sess)
		if err != nil {
			return Paged[api.Session]{}, err
		}

		outSessions = append(outSessions, res)
	}

	out := Paged[api.Session]{
		Data: outSessions,
	}

	return out, nil
}

func (idp *IdentityProvider) CreateRecoveryCodeForIdentity(ctx context.Context, request api.CreateRecoveryCodeForIdentityRequestObject) (*api.RecoveryCodeForIdentity, error) {
	req := idp.kratos.Admin.IdentityAPI.CreateRecoveryCodeForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryCodeForIdentityBody(orykratos.CreateRecoveryCodeForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
			FlowType:   request.Body.FlowType,
		})
	}

	recoveryCodeForIdentity, recoveryCodeForIdentityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = recoveryCodeForIdentityResp.Body.Close() }()

	resp, err := RecoveryCodeForIdentityToResponse(*recoveryCodeForIdentity)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (idp *IdentityProvider) CreateRecoveryLinkForIdentity(ctx context.Context, request api.CreateRecoveryLinkForIdentityRequestObject) (*api.RecoveryLinkForIdentity, error) {
	req := idp.kratos.Admin.IdentityAPI.CreateRecoveryLinkForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryLinkForIdentityBody(orykratos.CreateRecoveryLinkForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
		})
	}

	recoveryLinkForIdentity, recoveryLinkForIdentityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = recoveryLinkForIdentityResp.Body.Close() }()

	resp, err := RecoveryLinkForIdentityToResponse(*recoveryLinkForIdentity)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func createUserRelationTuples(ctx context.Context, c *keto.Client, identityID string) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customer",
					Relation:  "members",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    shared.AuthzIdentityID(identityID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    shared.AuthzIdentityID(identityID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert identity relation tuples: %w", err)
	}

	return nil
}

func deleteUserRelationTuples(ctx context.Context, c *keto.Client, identityID string) error {
	ownersResp, err := c.Read.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: new("Identity"),
			Object:    new(shared.AuthzIdentityID(identityID)),
			Relation:  new("owners"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list identity relation tuples: %w", err)
	}

	tuplesToDelete := []*rts.RelationTupleDelta{
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "customer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "developer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(identityID),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		},
	}

	for _, tuple := range ownersResp.RelationTuples {
		subject := tuple.GetSubject()
		if subject == nil {
			continue
		}

		tuplesToDelete = append(tuplesToDelete, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(identityID),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(subject.GetId()),
			},
		})
	}

	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuplesToDelete,
	}); err != nil {
		return fmt.Errorf("failed to delete identity relation tuples: %w", err)
	}

	return nil
}

type UserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

func (idp *IdentityProvider) GetUserInfo(ctx context.Context, userID uuid.UUID) (UserInfo, error) {
	// currently most of the stuff is in kratos traits and not my custom user table
	identity, err := idp.GetIdentity(ctx, api.GetIdentityRequestObject{ID: userID.String()})
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to fetch user info: %w", err)
	}

	traits, ok := identity.Traits.(map[string]any)
	if !ok {
		return UserInfo{}, fmt.Errorf("failed to parse user traits")
	}

	username, ok1 := traits["username"].(string)
	firstName, ok2 := traits["first_name"].(string)
	lastName, ok3 := traits["first_name"].(string)

	avatarURL, ok4 := traits["avatar_url"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return UserInfo{}, fmt.Errorf("failed to parse user traits values")
	}

	resp := UserInfo{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		AvatarURL: avatarURL,
	}

	return resp, nil
}

// GetUsername just gets username for game related stuff
func (idp *IdentityProvider) GetUsername(ctx context.Context, userID uuid.UUID) (string, error) {
	// currently most of the stuff is in kratos traits and not my custom user table
	identity, err := idp.GetIdentity(ctx, api.GetIdentityRequestObject{ID: userID.String()})
	if err != nil {
		return "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	traits, ok := identity.Traits.(map[string]any)
	if !ok {
		return "", fmt.Errorf("failed to parse user traits")
	}

	username, ok := traits["username"].(string)
	if !ok {
		return "", fmt.Errorf("no username found")
	}

	return username, nil
}

func (idp *IdentityProvider) CreateFriendRequest(ctx context.Context, request api.CreateFriendRequestRequestObject) (api.CreateFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API CreateFriendRequest")
}

func (idp *IdentityProvider) ListFriendRequests(ctx context.Context, request api.ListFriendRequestsRequestObject) (api.ListFriendRequestsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFriendRequests")
}

func (idp *IdentityProvider) AcceptFriendRequest(ctx context.Context, request api.AcceptFriendRequestRequestObject) (api.AcceptFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API AcceptFriendRequest")
}

func (idp *IdentityProvider) DeclineFriendRequest(ctx context.Context, request api.DeclineFriendRequestRequestObject) (api.DeclineFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeclineFriendRequest")
}

func (idp *IdentityProvider) CancelFriendRequest(ctx context.Context, request api.CancelFriendRequestRequestObject) (api.CancelFriendRequestResponseObject, error) {
	panic("@TODO: IMPLEMENT API CancelFriendRequest")
}

func (idp *IdentityProvider) ListFriends(ctx context.Context, request api.ListFriendsRequestObject) (api.ListFriendsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFriends")
}

func (idp *IdentityProvider) GetFriend(ctx context.Context, request api.GetFriendRequestObject) (api.GetFriendResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetFriend")
}

func (idp *IdentityProvider) DeleteFriend(ctx context.Context, request api.DeleteFriendRequestObject) (api.DeleteFriendResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeleteFriend")
}

func (idp *IdentityProvider) DeleteFriends(ctx context.Context, request api.DeleteFriendsRequestObject) (api.DeleteFriendsResponseObject, error) {
	panic("@TODO: IMPLEMENT API DeleteFriends")
}

func (idp *IdentityProvider) ListFollowings(ctx context.Context, request api.ListFollowingsRequestObject) (api.ListFollowingsResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListFollowings")
}

func (idp *IdentityProvider) GetFollowing(ctx context.Context, request api.GetFollowingRequestObject) (api.GetFollowingResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetFollowing")
}

func (idp *IdentityProvider) FollowUser(ctx context.Context, request api.FollowUserRequestObject) (api.FollowUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API FollowUser")
}

func (idp *IdentityProvider) UnfollowUser(ctx context.Context, request api.UnfollowUserRequestObject) (api.UnfollowUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnfollowUser")
}

func (idp *IdentityProvider) UnfollowUsers(ctx context.Context, request api.UnfollowUsersRequestObject) (api.UnfollowUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnfollowUsers")
}

func (idp *IdentityProvider) ListBlockedUsers(ctx context.Context, request api.ListBlockedUsersRequestObject) (api.ListBlockedUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API ListBlockedUsers")
}

func (idp *IdentityProvider) GetBlockedUser(ctx context.Context, request api.GetBlockedUserRequestObject) (api.GetBlockedUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API GetBlockedUser")
}

func (idp *IdentityProvider) BlockUser(ctx context.Context, request api.BlockUserRequestObject) (api.BlockUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API BlockUser")
}

func (idp *IdentityProvider) UnblockUser(ctx context.Context, request api.UnblockUserRequestObject) (api.UnblockUserResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnblockUser")
}

func (idp *IdentityProvider) UnblockUsers(ctx context.Context, request api.UnblockUsersRequestObject) (api.UnblockUsersResponseObject, error) {
	panic("@TODO: IMPLEMENT API UnblockUsers")
}
