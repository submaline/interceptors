package interceptors

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/submaline/interceptors/internal"
	greetv1 "github.com/submaline/interceptors/internal/gen/greet/v1"
	"github.com/submaline/interceptors/internal/gen/greet/v1/greetv1connect"
	"github.com/submaline/interceptors/internal/service"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestLaunchUpServer(t *testing.T) {
	t.Parallel()

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to set up firebase app: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		t.Fatalf("failed to set up auth client: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("failed to set up logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	greetHandler := &service.GreetService{}
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(
		NewFirebaseAuthInterceptor(authClient, AuthPolicy{"/greet.v1.GreetService/PlainGreet": false}),
		NewLoggingInterceptor(logger))
	mux.Handle(greetv1connect.NewGreetServiceHandler(
		greetHandler,
		interceptors))

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Service listening on %v", port)
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		t.Fatalf("failed to serve: %v", err)
	}
}

func TestGreetRequest(t *testing.T) {
	t.Parallel()
	greetClient := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://localhost:8080"))

	req := connect.NewRequest(&greetv1.GreetRequest{Name: "john"})

	// トークンなし
	_, err := greetClient.Greet(context.Background(), req)
	if err != nil {
		// 意図したエラー
		t.Logf("failed to greet(no auth): %v", err)
	}

	// トークンつける
	tokenData, err := internal.GenToken(os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	if err != nil {
		// 意図しないエラー
		t.Fatalf("failed to generate token: %v", err)
	}
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenData.IdToken))

	// トークンあり
	resp, err := greetClient.Greet(context.Background(), req)
	if err != nil {
		// 意図しないエラー
		t.Fatalf("failed to greet(auth): %v", err)
	}
	log.Println(resp.Msg)
}

func TestPlainGreetRequest(t *testing.T) {
	t.Parallel()
	greetClient := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://localhost:8080"))

	// トークンなし
	req := connect.NewRequest(&greetv1.PlainGreetRequest{Name: "tom"})
	resp, err := greetClient.PlainGreet(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to plaingreet: %v", err)
	}

	log.Println(resp.Msg)
}

func TestStreamGreetRequest(t *testing.T) {
	t.Parallel()
	greetClient := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://localhost:8080"))

	req := connect.NewRequest(&greetv1.StreamGreetRequest{Name: "sam"})

	// トークンなし
	stream, err := greetClient.StreamGreet(context.Background(), req)
	if err != nil {
		t.Logf("failed to start streamGreet: %v", err)
	}
	for stream.Receive() {
		msg := stream.Msg()
		log.Println(msg)
	}
	if err = stream.Err(); err != nil {
		t.Logf("failed to get msg from streamGreet: %v", err)
	}

	// トークンつける
	tokenData, err := internal.GenToken(os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	if err != nil {
		// 意図しないエラー
		t.Fatalf("failed to generate token: %v", err)
	}
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenData.IdToken))

	// トークンあり
	stream, err = greetClient.StreamGreet(context.Background(), req)
	if err != nil {
		// 意図しないエラー
		t.Errorf("failed to start streamGreet: %v", err)
	}
	for stream.Receive() {
		msg := stream.Msg()
		log.Println(msg)
	}
	if err = stream.Err(); err != nil {
		// 意図しないエラー
		t.Errorf("failed to get msg from streamGreet: %v", err)
	}
}
