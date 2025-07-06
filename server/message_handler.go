package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/dto"
	"github.com/ory/client-go"
)

func (h *ApiHandler) ListCourierMessages(ctx context.Context, request api.ListCourierMessagesRequestObject) (api.ListCourierMessagesResponseObject, error) {
	req := h.Kratos.Admin.CourierAPI.ListCourierMessages(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil {
		req = req.PageToken(*request.Params.PageToken)
	}
	if request.Params.Recipient != nil {
		req = req.Recipient(*request.Params.Recipient)
	}
	if request.Params.Status != nil {
		req = req.Status(client.CourierMessageStatus(*request.Params.Status))
	}
	courierMessages, _, err := req.Execute()
	if err != nil {
		return make(api.ListCourierMessages200JSONResponse, 0), nil
	}
	resp := make(api.ListCourierMessages200JSONResponse, 0)
	for _, message := range courierMessages {
		res, err := dto.MessageToResponse(message)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func (h *ApiHandler) GetCourierMessage(ctx context.Context, request api.GetCourierMessageRequestObject) (api.GetCourierMessageResponseObject, error) {
	req := h.Kratos.Admin.CourierAPI.GetCourierMessage(ctx, request.ID)
	courierMessage, _, err := req.Execute()
	if err != nil {
		return api.GetCourierMessage404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Courier message not found"}}, nil
	}
	resp, err := dto.MessageToResponse(*courierMessage)
	if err != nil {
		return nil, err
	}
	return api.GetCourierMessage200JSONResponse(resp), nil
}
