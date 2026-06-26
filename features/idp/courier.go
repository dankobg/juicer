package idp

import (
	"context"

	api "github.com/dankobg/juicer/api/gen"
	orykratos "github.com/ory/client-go"
)

func (idp *IdentityProvider) ListCourierMessages(ctx context.Context, request api.ListCourierMessagesRequestObject) (Paged[api.Message], error) {
	req := idp.kratos.Admin.CourierAPI.ListCourierMessages(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	if request.Params.Recipient != nil {
		req = req.Recipient(*request.Params.Recipient)
	}

	if request.Params.Status != nil {
		req = req.Status(orykratos.CourierMessageStatus(*request.Params.Status))
	}

	courierMessages, courierMessagesResp, err := req.Execute()
	if err != nil {
		return Paged[api.Message]{}, err
	}

	defer func() { _ = courierMessagesResp.Body.Close() }()

	outMsgs := make([]api.Message, len(courierMessages))
	for _, msg := range courierMessages {
		res, err := MessageToResponse(msg)
		if err != nil {
			return Paged[api.Message]{}, err
		}

		outMsgs = append(outMsgs, res)
	}

	out := Paged[api.Message]{
		Data: outMsgs,
	}

	return out, nil
}

func (idp *IdentityProvider) GetCourierMessage(ctx context.Context, request api.GetCourierMessageRequestObject) (*api.Message, error) {
	req := idp.kratos.Admin.CourierAPI.GetCourierMessage(ctx, request.ID)

	courierMessage, courierMessageResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = courierMessageResp.Body.Close() }()

	out, err := MessageToResponse(*courierMessage)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
