package middleware

import (
	"net/http"

	"awsems/internal/cognito"
	"awsems/internal/utils"
)

func Authentication(jwtVerifier *cognito.JWTVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("access_token")
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Access token not found")
				return
			}

			_, err = jwtVerifier.VerifyAccessToken(cookie.Value)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Access token invalid or expired")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
