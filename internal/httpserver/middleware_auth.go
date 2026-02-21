package httpserver

import (
	"net/http"
	"strings"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/handlers"
)

func AuthMiddleware(expectedToken string) func(http.Handler) http.Handler {
	// Se não tiver token configurado, por segurança bloqueia tudo (exceto /health, que tratamos no router)
	if expectedToken == "" {
		expectedToken = "__MISSING_API_TOKEN__"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				handlers.WriteError(w, r, "UNAUTHORIZED", "missing Authorization header", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				handlers.WriteError(w, r, "UNAUTHORIZED", "invalid Authorization format (expected Bearer token)", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
			if token != expectedToken {
				handlers.WriteError(w, r, "UNAUTHORIZED", "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
