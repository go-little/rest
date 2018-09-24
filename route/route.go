package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/handlers"

	"github.com/go-little/rest/middleware"
)

//MainHandler cria o handler principal e configura as rotas da API
func MainHandler(routes Routes) http.Handler {
	mainRouter := mux.NewRouter()

	authRoutes := make([]*mux.Route, 0)

	fmt.Printf("\n################ Routes ################\n")
	for _, route := range routes {
		r := mainRouter.HandleFunc(route.pathPattern, route.handlerFunc).Methods(route.method)
		if route.auth {
			authRoutes = append(authRoutes, r)
		}
		fmt.Printf("%s %s\n", route.method, route.pathPattern)
	}
	fmt.Printf("################ Routes ################\n\n")

	mainRouter.Use(middleware.TracerMiddleware)
	mainRouter.Use(middleware.JWTMiddleware(authRoutes))

	mainHandler := handlers.RecoveryHandler()(mainRouter)
	return mainHandler
}
