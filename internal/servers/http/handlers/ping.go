package handlers

import (
	"context"

	api "github.com/Onnywrite/ssonny/api/oapi"
)

type InternalHandler struct{}

func (*InternalHandler) GetPing(ctx context.Context,
	request api.GetPingRequestObject,
) (api.GetPingResponseObject, error) {
	return api.GetPing200TextResponse("pong"), nil
}

func (*InternalHandler) GetHealthz(ctx context.Context,
	request api.GetHealthzRequestObject,
) (api.GetHealthzResponseObject, error) {
	return api.GetHealthz200TextResponse("ok"), nil
}

func (*InternalHandler) GetMetrics(ctx context.Context,
	request api.GetMetricsRequestObject,
) (api.GetMetricsResponseObject, error) {
	return api.GetMetrics200JSONResponse{}, nil
}
