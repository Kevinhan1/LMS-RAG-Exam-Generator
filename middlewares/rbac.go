package middlewares

import (
	"net/http"
)

func RequireRole(roleID int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid, ok := r.Context().Value(CtxRoleID).(int64)
			if !ok || rid != roleID {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
