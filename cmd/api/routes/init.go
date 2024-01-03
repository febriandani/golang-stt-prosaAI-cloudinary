package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/handler/api"
	"github.com/sirupsen/logrus"
)

func GetCoreEndpoint(conf *general.AppService, handler api.Handler, log *logrus.Logger) *mux.Router {
	parentRoute := mux.NewRouter()

	jwtRoute := parentRoute.PathPrefix("").Subrouter()
	nonJWTRoute := parentRoute.PathPrefix("").Subrouter()
	publicRoute := parentRoute.PathPrefix("").Subrouter()
	freeRoute := parentRoute.PathPrefix("").Subrouter()
	wsRoute := parentRoute.PathPrefix("").Subrouter()

	// Renew Access Token Endpoint.
	publicRoute.HandleFunc("/v1/token/renew", handler.Token.RenewAccessToken).Methods(http.MethodGet)

	// Middleware for public API
	nonJWTRoute.Use(handler.Public.AuthValidator)

	// Middleware
	if conf.Authorization.JWT.IsActive {
		log.Info("JWT token is active")
		jwtRoute.Use(handler.Token.JWTValidator)
	}

	// Get Endpoint.
	getV1(freeRoute, nonJWTRoute, jwtRoute, wsRoute, conf, handler)

	return parentRoute
}
