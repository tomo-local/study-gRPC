package middleware

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	// AuthServiceのJWT管理機能をインポート
	"auth/auth" // これは共有ライブラリとして利用
)

// AuthInterceptor はJWT認証を行うgRPCインターセプター
type AuthInterceptor struct {
	jwtManager *auth.JWTManager
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(jwtManager *auth.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}

// Unary は単項RPCのための認証インターセプター
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("Unary interceptor for method: %s", info.FullMethod)

		// 認証が不要なメソッドはスキップ（必要に応じて設定）
		if interceptor.isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// authorize はJWTトークンを検証する
func (interceptor *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	token := authHeader[7:] // "Bearer " を除去
	claims, err := interceptor.jwtManager.VerifyToken(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	// コンテキストにユーザー情報を追加（後で使用可能）
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.Email)
	ctx = context.WithValue(ctx, "user_name", claims.Name)

	return nil
}

// isPublicMethod は認証不要なメソッドかどうかを判定
func (interceptor *AuthInterceptor) isPublicMethod(method string) bool {
	// 認証不要なメソッドのリスト（必要に応じて拡張）
	publicMethods := []string{
		"/note.NoteService/HealthCheck", // 例：ヘルスチェック
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}

	return false
}
