// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: greet/v1/greet.proto

package greetv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/submaline/interceptors/internal/gen/greet/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// GreetServiceName is the fully-qualified name of the GreetService service.
	GreetServiceName = "greet.v1.GreetService"
)

// GreetServiceClient is a client for the greet.v1.GreetService service.
type GreetServiceClient interface {
	Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error)
	PlainGreet(context.Context, *connect_go.Request[v1.PlainGreetRequest]) (*connect_go.Response[v1.PlainGreetResponse], error)
	StreamGreet(context.Context, *connect_go.Request[v1.StreamGreetRequest]) (*connect_go.ServerStreamForClient[v1.StreamGreetResponse], error)
}

// NewGreetServiceClient constructs a client for the greet.v1.GreetService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGreetServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) GreetServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &greetServiceClient{
		greet: connect_go.NewClient[v1.GreetRequest, v1.GreetResponse](
			httpClient,
			baseURL+"/greet.v1.GreetService/Greet",
			opts...,
		),
		plainGreet: connect_go.NewClient[v1.PlainGreetRequest, v1.PlainGreetResponse](
			httpClient,
			baseURL+"/greet.v1.GreetService/PlainGreet",
			opts...,
		),
		streamGreet: connect_go.NewClient[v1.StreamGreetRequest, v1.StreamGreetResponse](
			httpClient,
			baseURL+"/greet.v1.GreetService/StreamGreet",
			opts...,
		),
	}
}

// greetServiceClient implements GreetServiceClient.
type greetServiceClient struct {
	greet       *connect_go.Client[v1.GreetRequest, v1.GreetResponse]
	plainGreet  *connect_go.Client[v1.PlainGreetRequest, v1.PlainGreetResponse]
	streamGreet *connect_go.Client[v1.StreamGreetRequest, v1.StreamGreetResponse]
}

// Greet calls greet.v1.GreetService.Greet.
func (c *greetServiceClient) Greet(ctx context.Context, req *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error) {
	return c.greet.CallUnary(ctx, req)
}

// PlainGreet calls greet.v1.GreetService.PlainGreet.
func (c *greetServiceClient) PlainGreet(ctx context.Context, req *connect_go.Request[v1.PlainGreetRequest]) (*connect_go.Response[v1.PlainGreetResponse], error) {
	return c.plainGreet.CallUnary(ctx, req)
}

// StreamGreet calls greet.v1.GreetService.StreamGreet.
func (c *greetServiceClient) StreamGreet(ctx context.Context, req *connect_go.Request[v1.StreamGreetRequest]) (*connect_go.ServerStreamForClient[v1.StreamGreetResponse], error) {
	return c.streamGreet.CallServerStream(ctx, req)
}

// GreetServiceHandler is an implementation of the greet.v1.GreetService service.
type GreetServiceHandler interface {
	Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error)
	PlainGreet(context.Context, *connect_go.Request[v1.PlainGreetRequest]) (*connect_go.Response[v1.PlainGreetResponse], error)
	StreamGreet(context.Context, *connect_go.Request[v1.StreamGreetRequest], *connect_go.ServerStream[v1.StreamGreetResponse]) error
}

// NewGreetServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGreetServiceHandler(svc GreetServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/greet.v1.GreetService/Greet", connect_go.NewUnaryHandler(
		"/greet.v1.GreetService/Greet",
		svc.Greet,
		opts...,
	))
	mux.Handle("/greet.v1.GreetService/PlainGreet", connect_go.NewUnaryHandler(
		"/greet.v1.GreetService/PlainGreet",
		svc.PlainGreet,
		opts...,
	))
	mux.Handle("/greet.v1.GreetService/StreamGreet", connect_go.NewServerStreamHandler(
		"/greet.v1.GreetService/StreamGreet",
		svc.StreamGreet,
		opts...,
	))
	return "/greet.v1.GreetService/", mux
}

// UnimplementedGreetServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGreetServiceHandler struct{}

func (UnimplementedGreetServiceHandler) Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("greet.v1.GreetService.Greet is not implemented"))
}

func (UnimplementedGreetServiceHandler) PlainGreet(context.Context, *connect_go.Request[v1.PlainGreetRequest]) (*connect_go.Response[v1.PlainGreetResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("greet.v1.GreetService.PlainGreet is not implemented"))
}

func (UnimplementedGreetServiceHandler) StreamGreet(context.Context, *connect_go.Request[v1.StreamGreetRequest], *connect_go.ServerStream[v1.StreamGreetResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("greet.v1.GreetService.StreamGreet is not implemented"))
}