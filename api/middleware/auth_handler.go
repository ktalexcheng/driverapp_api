package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ktalexcheng/trailbrake_api/api/model"
	"github.com/ktalexcheng/trailbrake_api/util"
)

// Validate token of requests
func AuthHandler(mg *util.MongoClient) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return ErrorHandler(
			func(w http.ResponseWriter, r *http.Request) error {
				// Get the authorization header value
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Missing authorization token"))
					return nil
				}

				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				token := model.Token{
					TokenString: tokenString,
				}

				// Verify token is valid and get user information from database
				err := token.VerifyToken(mg)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return err
				}

				ctx := context.WithValue(r.Context(), model.TokenKey, token)

				next.ServeHTTP(w, r.WithContext(ctx))
				return nil
			},
		)
	}
}
