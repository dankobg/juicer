package idp

import (
	"context"

	api "github.com/dankobg/juicer/api/gen"
)

func (idp *IdentityProvider) ListSessions(ctx context.Context, request api.ListSessionsRequestObject) (Paged[api.Session], error) {
	req := idp.kratos.Admin.IdentityAPI.ListSessions(ctx)
	if request.Params.Active != nil {
		req = req.Active(*request.Params.Active)
	}

	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	if request.Params.Expand != nil {
		expands := make([]string, 0, len(*request.Params.Expand))
		for _, x := range *request.Params.Expand {
			expands = append(expands, string(x))
		}

		req = req.Expand(expands)
	}

	sessions, sessionsResp, err := req.Execute()
	if err != nil {
		return Paged[api.Session]{}, err
	}

	defer func() { _ = sessionsResp.Body.Close() }()

	outSessions := make([]api.Session, len(sessions))
	for _, session := range sessions {
		res, err := SessionToResponse(session)
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

func (idp *IdentityProvider) GetSession(ctx context.Context, request api.GetSessionRequestObject) (*api.Session, error) {
	req := idp.kratos.Admin.IdentityAPI.GetSession(ctx, request.ID)
	if request.Params.Expand != nil {
		expands := make([]string, 0, len(*request.Params.Expand))
		for _, x := range *request.Params.Expand {
			expands = append(expands, string(x))
		}

		req = req.Expand(expands)
	}

	session, sessionResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = sessionResp.Body.Close() }()

	out, err := SessionToResponse(*session)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (idp *IdentityProvider) DisableSession(ctx context.Context, request api.DisableSessionRequestObject) error {
	req := idp.kratos.Admin.IdentityAPI.DisableSession(ctx, request.ID)

	disableSessionResp, err := req.Execute()
	if err != nil {
		return err
	}

	defer func() { _ = disableSessionResp.Body.Close() }()

	return nil
}

func (idp *IdentityProvider) ExtendSession(ctx context.Context, request api.ExtendSessionRequestObject) (*api.Session, error) {
	req := idp.kratos.Admin.IdentityAPI.ExtendSession(ctx, request.ID)

	session, sessionResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = sessionResp.Body.Close() }()

	out, err := SessionToResponse(*session)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
