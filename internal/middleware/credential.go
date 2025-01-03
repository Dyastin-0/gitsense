package middleware

import (
	"net/http"

	"github.com/Dyastin-0/gitsense/internal/config"
)

func Credential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fetchSite := r.Header.Get("Sec-Fetch-Site")
		isAllowed := true

		if fetchSite == "same-origin" || fetchSite == "same-site" {
			isAllowed = true
		} else if origin != "" {
			for _, allowedOrigin := range config.AllowedOrigins {
				if origin == allowedOrigin {
					isAllowed = true
					break
				}
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		} else {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
