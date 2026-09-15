package middleware

import (
	"net/http"

	"awsems/internal/cognito"
	"awsems/internal/utils"
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

			_, err := jwtVerifier.VerifyAccessToken(tokenString)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Access token invalid or expired")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
