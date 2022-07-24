package interceptors

import (
	"context"
	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type loggingInterceptor struct {
	logger *zap.Logger
}

func NewLoggingInterceptor(logger *zap.Logger) *loggingInterceptor {
	return &loggingInterceptor{
		logger: logger,
	}
}

func (l *loggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		resp, err := next(ctx, request)
		if err != nil {
			l.logger.Error(
				"an error occurred",
				zap.String("method", request.Spec().Procedure),
				zap.Errors("detail", []error{err}),
			)
		} else {
			l.logger.Info(
				"finish",
				zap.String("method", request.Spec().Procedure),
			)
		}
		return resp, err
	})
}

// WrapStreamingClient maybe client stream
func (l *loggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return &logHeaderInspectingClientConn{
			StreamingClientConn: next(ctx, spec),
			logger:              l.logger,
		}
	}
}

type logHeaderInspectingClientConn struct {
	connect.StreamingClientConn
	logger *zap.Logger
}

func (cc *logHeaderInspectingClientConn) Send(msg any) error {
	err := cc.StreamingClientConn.Send(msg)
	if err != nil {
		cc.logger.Error(
			"an error occurred",
			zap.String("method", cc.Spec().Procedure),
			zap.Errors("detail", []error{err}),
		)
	} else {
		cc.logger.Info(
			"finish",
			zap.String("method", cc.Spec().Procedure),
		)
	}
	return err
}
func (cc *logHeaderInspectingClientConn) Receive(msg any) error {
	err := cc.StreamingClientConn.Receive(msg)
	if err != nil {
		cc.logger.Error(
			"an error occurred",
			zap.String("method", cc.Spec().Procedure),
			zap.Errors("detail", []error{err}),
		)
	} else {
		cc.logger.Info(
			"finish",
			zap.String("method", cc.Spec().Procedure),
		)
	}
	return err
}

type logHeaderInspectingHandlerConn struct {
	connect.StreamingHandlerConn
	logger *zap.Logger
}

func (hc *logHeaderInspectingHandlerConn) Send(msg any) error {
	err := hc.StreamingHandlerConn.Send(msg)
	if err != nil {
		hc.logger.Error(
			"an error occurred",
			zap.String("method", hc.Spec().Procedure),
			zap.Errors("detail", []error{err}),
		)
	} else {
		hc.logger.Info(
			"finish",
			zap.String("method", hc.Spec().Procedure),
		)
	}
	return err
}
func (hc *logHeaderInspectingHandlerConn) Receive(msg any) error {
	err := hc.StreamingHandlerConn.Receive(msg)
	if err != nil {
		hc.logger.Error(
			"an error occurred",
			zap.String("method", hc.Spec().Procedure),
			zap.Errors("detail", []error{err}),
		)
	} else {
		hc.logger.Info(
			"finish",
			zap.String("method", hc.Spec().Procedure),
		)
	}
	return err
}

// WrapStreamingHandler maybe server side stream
func (l *loggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(ctx, &logHeaderInspectingHandlerConn{
			conn,
			l.logger,
		})
	}
}
