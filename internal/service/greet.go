package service

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	greetv1 "github.com/submaline/interceptors/internal/gen/greet/v1"
	"log"
	"strconv"
)

const (
	applicationName = "Submaline"
)

type GreetService struct{}

func (g GreetService) Greet(_ context.Context, request *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	userId := request.Header().Get(fmt.Sprintf("X-%s-UserId", applicationName))
	isAdmin_ := request.Header().Get(fmt.Sprintf("X-%s-Admin", applicationName))
	isAdmin, err := strconv.ParseBool(isAdmin_)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to parse isAdmin: %v", err))
	}
	kind := "normal"
	if isAdmin {
		kind = "admin"
	}

	return connect.NewResponse(
		&greetv1.GreetResponse{
			Greeting: fmt.Sprintf("hello, %s(%s: %s)!", request.Msg.Name, kind, userId),
		}), nil
}

func (g GreetService) PlainGreet(_ context.Context, request *connect.Request[greetv1.PlainGreetRequest]) (*connect.Response[greetv1.PlainGreetResponse], error) {
	return connect.NewResponse(
		&greetv1.PlainGreetResponse{
			Greeting: fmt.Sprintf("hello, %s!", request.Msg.Name)}), nil
}

func (g GreetService) StreamGreet(ctx context.Context, request *connect.Request[greetv1.StreamGreetRequest], stream *connect.ServerStream[greetv1.StreamGreetResponse]) error {
	log.Println("StreamGreet REQ", request.Header())
	userId := request.Header().Get(fmt.Sprintf("X-%s-UserId", applicationName))
	isAdmin_ := request.Header().Get(fmt.Sprintf("X-%s-Admin", applicationName))
	isAdmin, err := strconv.ParseBool(isAdmin_)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to parse isAdmin: %v", err))
	}
	kind := "normal"
	if isAdmin {
		kind = "admin"
	}

	count := 0
streamLoop:
	for {
		select {
		case <-ctx.Done():
			// stream closed, ex. client disconnected
			break streamLoop
		default:
			err := stream.Send(&greetv1.StreamGreetResponse{
				Greeting: fmt.Sprintf("(%d) hello, %s(%s: %s)!", count, request.Msg.Name, kind, userId)})
			if err != nil {
				return err
			}
			count++
			if count >= 10 {
				break streamLoop
			}
		}
	}

	return ctx.Err()
}
