package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RequireRoles(roles ...int) mux.MiddlewareFunc {
	roleSet := make(map[int]bool)
	for _, r := range roles {
		roleSet[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleID, ok := r.Context().Value(CtxRoleID).(int)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !roleSet[roleID] {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
