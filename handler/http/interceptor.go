package http

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	h "github.com/yoshino-s/go-framework/handlers/http"
)

var _ connect.Interceptor = (*interceptor)(nil)

type interceptor struct {
}

func (i *interceptor) handleError(ctx echo.Context, err error) {
	hub := sentryecho.GetHubFromContext(ctx)
	if hub == nil {
		fmt.Println("hub is nil")
		return
	}

	r := ctx.Request()
	hub.RecoverWithContext(
		context.WithValue(r.Context(), sentry.RequestContextKey, r),
		err,
	)
}

func (i *interceptor) WrapUnary(f connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		r, err := f(ctx, ar)
		c, ok := h.EchoContextFromContext(ctx)
		if ok && err != nil {
			i.handleError(c, err)
		}

		return r, err
	})
}

func (i *interceptor) WrapStreamingHandler(f connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		err := f(ctx, shc)
		c, ok := h.EchoContextFromContext(ctx)
		if ok && err != nil {
			i.handleError(c, err)
		}

		return err
	})
}

func (i *interceptor) WrapStreamingClient(f connect.StreamingClientFunc) connect.StreamingClientFunc {
	return f
}
