package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	handlerPort "kiramishima/m-backend/internal/core/ports/handlers"
	svcports "kiramishima/m-backend/internal/core/ports/services"
	"net/http"
	"os"
)

var _ handlerPort.UserHandlers = (*UserHandlers)(nil)

// NewBondHandlers creates a instance of auth handlers
func NewUserHandlers(r *chi.Mux, logger *zap.SugaredLogger, s svcports.UserService, render *render.Render, validate *validator.Validate) {
	var tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SecretKey")), nil)

	handler := &UserHandlers{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/me", func(r chi.Router) {
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Get("/", handler.GetProfileHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/", handler.UpdateProfileHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Put("/bonds", handler.GetUserBondsHandler)
	})
}

type UserHandlers struct {
	logger   *zap.SugaredLogger
	service  svcports.UserService
	response *render.Render
	validate *validator.Validate
}

func (u UserHandlers) GetProfileHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandlers) UpdateProfileHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandlers) GetUserBondsHandler(w http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}
