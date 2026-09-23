package middleware

import (
	"context"
	"net/http"

	"awsems/internal/cognito"
	"awsems/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

func Authentication(jwtVerifier *cognito.JWTVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Error(w, http.StatusUnauthorized, "Access token not found")
				return
			}

			// Extract the token assuming format: Bearer <token>
			tokenString := ""
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				tokenString = authHeader
			}

			token, err := jwtVerifier.VerifyAccessToken(tokenString)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Access token invalid or expired")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.Error(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				utils.Error(w, http.StatusUnauthorized, "sub missing from access token")
				return
			}

			// Add the sub to the context
			ctx := context.WithValue(r.Context(), "user_sub", sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
