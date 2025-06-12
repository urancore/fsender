package middleware

import (
	"fsender/internal/utils"
	"fsender/pkg"

	"net/http"
	"fmt"
)

func (md *Middlewares) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_ip, err := utils.GetClientIP(r)
		if err != nil {
			user_ip = "not found"
			fmt.Println(err.Error())
		}

		next.ServeHTTP(w, r)
		pkg.Log(user_ip, r.Method, r.URL.Path)
	})
}
