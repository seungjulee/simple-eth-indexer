package server

import (
	"context"
	"net/http"

	"github.com/seungjulee/simple-eth-indexer/pkg/logger"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

// NewLoggingClientHooks logs request and errors to stdout in the client
func NewLoggingClientHooks() *twirp.ClientHooks {
    return &twirp.ClientHooks{
        RequestPrepared: func(ctx context.Context, r *http.Request) (context.Context, error) {
			logger.Info("Req: ", zap.String("host", r.Host), zap.String("url_path", r.URL.Path))
            return ctx, nil
        },
        Error: func(ctx context.Context, twerr twirp.Error) {
			logger.Info("Error: ", zap.String("err", string(twerr.Code())))
        },
        ResponseReceived: func(ctx context.Context) {
			logger.Info("success resp")
        },
    }
}
// NewLoggingClientHooks logs request and errors to stdout in the client
func NewLoggingServerHooks() *twirp.ServerHooks {
    return &twirp.ServerHooks{
        RequestRouted: func(ctx context.Context) (context.Context, error) {
            method, _ := twirp.MethodName(ctx)
            // log.Println("Method: " + method)
            logger.Info("method: ", zap.String("method", method))
            return ctx, nil
        },
        Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			logger.Info("Error: ", zap.String("err", string(twerr.Code())))
            return ctx
        },
        ResponseSent: func(ctx context.Context) {
            logger.Info("Response Sent")
        },
    }
}