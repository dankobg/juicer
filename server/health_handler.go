package server

import (
	"context"

	api "github.com/dankobg/juicer/api/gen"
)

func (h *ApiHandler) GetHealthAlive(ctx context.Context, request api.GetHealthAliveRequestObject) (api.GetHealthAliveResponseObject, error) {
	return api.GetHealthAlive200JSONResponse{Alive: true}, nil
}

func (h *ApiHandler) GetHealthReady(ctx context.Context, request api.GetHealthReadyRequestObject) (api.GetHealthReadyResponseObject, error) {
	return api.GetHealthReady200JSONResponse{Ready: true}, nil
}
