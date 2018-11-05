package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-little/rest/reply"
	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTMiddlewareConfig struct {
	JWTKey                string
	AuthorizationHeader   string
	UnathorizedStatusCode int
	UnathorizedBody       []byte
	UnathorizedJSON       interface{}
}

// JWTMiddleware comment
func JWTMiddleware(config JWTMiddlewareConfig, authRoutes []*mux.Route) func(next http.Handler) http.Handler {

	jwtKey := config.JWTKey
	if jwtKey != "" {
		fmt.Printf("JWT auth enabled with key: %s\n", jwtKey)
	}

	authorizationHeader := "authorization"
	if config.AuthorizationHeader != "" {
		authorizationHeader = config.AuthorizationHeader
	}

	unathorizedStatusCode := 401
	if config.UnathorizedStatusCode != 0 {
		unathorizedStatusCode = config.UnathorizedStatusCode
	}

	unathorizedBody := config.UnathorizedBody
	unathorizedJSON := config.UnathorizedJSON
	if unathorizedBody == nil && unathorizedJSON == nil {
		unathorizedJSON = map[string]string{"message": "unathorized"}
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

			if config.JWTKey != "" && authEnable {

				authorizationToken := r.Header.Get(authorizationHeader)
				token, err := jwt.Parse(authorizationToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(config.JWTKey), nil
				})

				if err != nil || !token.Valid {
					re := reply.StatusCode(unathorizedStatusCode)
					if unathorizedBody != nil {
						re = re.Body(unathorizedBody)
					} else if unathorizedJSON != nil {
						re = re.JSON(unathorizedJSON)
					}
					re.Do(w)
					return
				}

			}

			next.ServeHTTP(w, r)

		})
	}
}
