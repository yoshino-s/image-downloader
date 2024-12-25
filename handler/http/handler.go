package http

import (
	"context"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/handlers/http"
	"github.com/yoshino-s/go-framework/telemetry"
	"github.com/yoshino-s/image-downloader/gen"
	"github.com/yoshino-s/image-downloader/gen/v1/v1connect"
	"go.akshayshah.org/connectproto"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ application.Application = (*Handler)(nil)

type Handler struct {
	*http.Handler

	imageHandler v1connect.ImageServiceHandler
}

func New() *Handler {
	return &Handler{
		Handler: http.New(),
	}
}

func (h *Handler) SetImageHandler(handler v1connect.ImageServiceHandler) {
	h.imageHandler = handler
}

func (h *Handler) Setup(ctx context.Context) {
	common.MustNoNil(h.imageHandler)
	h.Handler.Setup(ctx)

	h.Swagger("/swagger", gen.OpenAPI)

	reflector := grpcreflect.NewStaticReflector(
		v1connect.ImageServiceName,
	)

	h.HandleGrpc(grpcreflect.NewHandlerV1(reflector))
	h.HandleGrpc(grpcreflect.NewHandlerV1Alpha(reflector))

	var opts []connect.HandlerOption

	opts = append(opts, connectproto.WithJSON(
		protojson.MarshalOptions{EmitUnpopulated: true, EmitDefaultValues: true},
		protojson.UnmarshalOptions{DiscardUnknown: true},
	))

	if telemetry.IsSentryInitialized() {
		opts = append(opts, connect.WithInterceptors(&interceptor{}))
	}

	h.HandleGrpc(v1connect.NewImageServiceHandler(h.imageHandler, opts...))
}

func (h *Handler) Run(ctx context.Context) {
	h.Handler.Run(ctx)
}
