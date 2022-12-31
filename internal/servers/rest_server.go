package servers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ProfileController interface {
	SaveProfile(http.ResponseWriter, *http.Request)
	GetProfiles(http.ResponseWriter, *http.Request)
}

type AuthController interface {
}

type RESTServer struct {
	server            http.Server
	authController    AuthController
	profileController ProfileController
	logger            *zap.Logger
}

func NewRESTServer(address string, authController AuthController, profileController ProfileController, logger *zap.Logger) RESTServer {
	handler := mux.NewRouter().StrictSlash(true)
	handler.HandleFunc("/profiles", profileController.GetProfiles).Methods("GET")
	handler.HandleFunc("/profile", profileController.SaveProfile).Methods("POST")

	return RESTServer{
		server:            http.Server{Addr: address, Handler: handler},
		authController:    authController,
		profileController: profileController,
		logger:            logger,
	}
}

func (s *RESTServer) ListenAndServe() {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("failed to start REST server" + err.Error())
	}
}

func (s *RESTServer) Shutdown(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		s.logger.Error("failed to close REST server" + err.Error())
	}
}
