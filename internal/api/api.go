package api

import (
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "wallet-service/docs"
	"wallet-service/internal/service"
)

type API struct {
	Wallets service.Wallets
}

func New() *API {
	return &API{
		Wallets: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r = r.PathPrefix("/api/v1/wallets").Subrouter()

	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Save), constants.ADMIN)).Methods(http.MethodPost)
	r.Handle("/{username}", middlewares.AuthMiddleware(http.HandlerFunc(a.FindByUser), constants.MODERATOR, constants.ADMIN)).Methods(http.MethodGet)
	r.Handle("/{username}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), constants.ADMIN)).Methods(http.MethodDelete)

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}
