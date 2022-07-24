package interceptors

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bufbuild/connect-go"
	"strings"
)

const (
	applicationName = "Submaline"
)

type firebaseAuthInterceptor struct {
	client *auth.Client
	policy AuthPolicy
}

func NewFirebaseAuthInterceptor(client *auth.Client, policy AuthPolicy) *firebaseAuthInterceptor {
	return &firebaseAuthInterceptor{
		client: client,
		policy: policy,
	}
}

func (i *firebaseAuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		// 関数のフルパスを取得
		funcFullPath := request.Spec().Procedure
		// ポリシーに認証は必要ないとされていたら本来の処理へ
		requireAuthorization, ok := i.policy[funcFullPath]
		if ok && !requireAuthorization {
			return next(ctx, request)
		}
		// ヘッダーからトークンを取り出してあげる
		// `Bearer ${TOKEN}`
		bearerIdToken := request.Header().Get("Authorization")
		// 入っていなかったらアウト
		if bearerIdToken == "" {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("%s requires authentication", funcFullPath))
		}

		// "Bearer "を消してidToken本体だけにしてあげる
		// `${TOKEN}`
		idToken := strings.Replace(bearerIdToken, "Bearer ", "", 1)

		// firebase sdkを使ってトークンが有効かを確認
		token, err := i.client.VerifyIDToken(context.Background(), idToken)
		// todo : エラーメッセージを作ってあげる
		if err != nil {
			switch {
			case auth.IsIDTokenExpired(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			case auth.IsIDTokenInvalid(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			case auth.IsSessionCookieInvalid(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			case auth.IsIDTokenRevoked(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			case auth.IsSessionCookieRevoked(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			case auth.IsUserNotFound(err):
				return nil, connect.NewError(
					connect.CodeUnauthenticated, err)
			default:
				return nil, connect.NewError(
					connect.CodeInternal, err)
			}
		}
		// ユーザーidをつけてあげる
		request.Header().Set(fmt.Sprintf("X-%s-UserId", applicationName), token.UID)

		// カスタムクレームを確認
		claims := token.Claims

		// アドミンであるかを確認
		admin, ok := claims["admin"]
		if ok && admin.(bool) {
			request.Header().Set(fmt.Sprintf("X-%s-Admin", applicationName), "true")
		} else {
			request.Header().Set(fmt.Sprintf("X-%s-Admin", applicationName), "false")
		}

		// 本来の処理へ
		return next(ctx, request)
	})
}

func (i *firebaseAuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		//log.Println(conn.RequestHeader())
		//log.Println(conn.ResponseHeader())
		return &headerInspectingClientConn{
			StreamingClientConn: next(ctx, spec),
			client:              i.client,
		}
	}
}

type headerInspectingClientConn struct {
	connect.StreamingClientConn
	client *auth.Client
	policy AuthPolicy
}

func (cc *headerInspectingClientConn) Send(msg any) error {
	//log.Println("headerInspectingClientConn.Send REQ", cc.RequestHeader())
	//log.Println("headerInspectingClientConn.Send RESP", cc.ResponseHeader())
	return cc.StreamingClientConn.Send(msg)
}
func (cc *headerInspectingClientConn) Receive(msg any) error {
	//log.Println("headerInspectingClientConn.Receive REQ", cc.RequestHeader())
	//log.Println("headerInspectingClientConn.Receive RESP", cc.ResponseHeader())
	return cc.StreamingClientConn.Receive(msg)
}

func (i *firebaseAuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		//log.Println("WrapStreamingHandler REQ", conn.RequestHeader())
		//log.Println("WrapStreamingHandler RESP", conn.ResponseHeader())
		return next(ctx, &headerInspectingHandlerConn{
			conn,
			i.client,
			i.policy,
		})
	}
}

type headerInspectingHandlerConn struct {
	connect.StreamingHandlerConn
	client *auth.Client
	policy AuthPolicy
}

func (hc *headerInspectingHandlerConn) Send(msg any) error {
	//hc.RequestHeader().Set("add", "hc req send")
	//hc.ResponseHeader().Set("add", "hc resp send")
	return hc.StreamingHandlerConn.Send(msg)
}

func (hc *headerInspectingHandlerConn) Receive(msg any) error {
	/// 関数のフルパスを取得
	funcFullPath := hc.Spec().Procedure
	// ポリシーに認証は必要ないとされていたら...
	requireAuthorization, ok := hc.policy[funcFullPath]
	if ok && !requireAuthorization {
		return nil
	}

	// ヘッダーからトークンを取り出してあげる
	// `Bearer ${TOKEN}`
	bearerIdToken := hc.RequestHeader().Get("Authorization")
	// 入っていなかったらアウト
	if bearerIdToken == "" {
		return connect.NewError(
			connect.CodeUnauthenticated,
			fmt.Errorf("%s requires authentication", funcFullPath))
	}

	// "Bearer "を消してidToken本体だけにしてあげる
	// `${TOKEN}`
	idToken := strings.Replace(bearerIdToken, "Bearer ", "", 1)

	// firebase sdkを使ってトークンが有効かを確認
	token, err := hc.client.VerifyIDToken(context.Background(), idToken)
	// todo : エラーメッセージを作ってあげる
	if err != nil {
		switch {
		case auth.IsIDTokenExpired(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		case auth.IsIDTokenInvalid(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		case auth.IsSessionCookieInvalid(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		case auth.IsIDTokenRevoked(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		case auth.IsSessionCookieRevoked(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		case auth.IsUserNotFound(err):
			return connect.NewError(
				connect.CodeUnauthenticated, err)
		default:
			return connect.NewError(
				connect.CodeInternal, err)
		}
	}
	// ユーザーidをつけてあげる
	hc.RequestHeader().Set(fmt.Sprintf("X-%s-UserId", applicationName), token.UID)

	// カスタムクレームを確認
	claims := token.Claims

	// アドミンであるかを確認
	admin, ok := claims["admin"]
	if ok && admin.(bool) {
		hc.RequestHeader().Set(fmt.Sprintf("X-%s-Admin", applicationName), "true")
	} else {
		hc.RequestHeader().Set(fmt.Sprintf("X-%s-Admin", applicationName), "false")
	}

	return hc.StreamingHandlerConn.Receive(msg)
}
