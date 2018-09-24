package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-little/rest/reply"
	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWTMiddleware comment
func JWTMiddleware(authRoutes []*mux.Route) func(next http.Handler) http.Handler {

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey != "" {
		fmt.Printf("JWT auth enabled with key: %v\n", jwtKey)
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			route := mux.CurrentRoute(r)
			authEnable := false

			for _, authRoute := range authRoutes {
				if authRoute == route {
					authEnable = true
					break
				}
			}

			if jwtKey != "" && authEnable {

				authorizationToken := r.Header.Get("authorization")
				token, err := jwt.Parse(authorizationToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(jwtKey), nil
				})

				if err != nil || !token.Valid {
					reply.StatusCode(401).JSON(map[string]string{"message": "unathorized"}).Do(w)
					return
				}

			}

			next.ServeHTTP(w, r)

		})
	}
}
