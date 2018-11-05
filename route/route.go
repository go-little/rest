package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/handlers"

	"github.com/go-little/rest/middleware"

	"github.com/go-little/rest/reply"
)

type MainHandlerConfig struct {
	NotFoundStatusCode int
	NotFoundBody       []byte
	NotFoundJSON       interface{}

	MethodNotAllowedStatusCode int
	MethodNotAllowedBody       []byte
	MethodNotAllowedJSON       interface{}

	TracerMiddlewareConfig middleware.TracerMiddlewareConfig
	JWTMiddlewareConfig    middleware.JWTMiddlewareConfig
}

//MainHandler cria o handler principal e configura as rotas da API
func MainHandler(config MainHandlerConfig, routes Routes) http.Handler {
	mainRouter := mux.NewRouter()

	// Resposta padrao em caso de nao encontrar a rota
	notFoundStatusCode := 404
	if config.NotFoundStatusCode != 0 {
		notFoundStatusCode = config.NotFoundStatusCode
	}

	notFoundBody := config.NotFoundBody
	notFoundJSON := config.NotFoundJSON
	if notFoundBody == nil && notFoundJSON == nil {
		notFoundJSON = map[string]string{"message": "not found"}
	}

	mainRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := reply.StatusCode(notFoundStatusCode)
		if notFoundBody != nil {
			re = re.Body(notFoundBody)
		} else if notFoundJSON != nil {
			re = re.JSON(notFoundJSON)
		}
		re.Do(w)
	})

	// Resposta padrao em caso de nao encontrar a rota e o metodo
	methodNotAllowedStatusCode := 404
	if config.MethodNotAllowedStatusCode != 0 {
		methodNotAllowedStatusCode = config.MethodNotAllowedStatusCode
	}

	methodNotAllowedBody := config.MethodNotAllowedBody
	methodNotAllowedJSON := config.MethodNotAllowedJSON
	if methodNotAllowedBody == nil && methodNotAllowedJSON == nil {
		notFoundJSON = map[string]string{"message": "method not allowed"}
	}

	mainRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := reply.StatusCode(methodNotAllowedStatusCode)
		if methodNotAllowedBody != nil {
			re = re.Body(methodNotAllowedBody)
		} else if methodNotAllowedJSON != nil {
			re = re.JSON(methodNotAllowedJSON)
		}
		re.Do(w)
	})

	// Registra todas as rotas
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

	// Registra todos os middlewares
	mainRouter.Use(middleware.TracerMiddleware(config.TracerMiddlewareConfig))
	mainRouter.Use(middleware.JWTMiddleware(config.JWTMiddlewareConfig, authRoutes))

	mainHandler := handlers.RecoveryHandler()(mainRouter)
	return mainHandler
}
