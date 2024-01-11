package bootstrap

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"kiramishima/m-backend/config"
	"kiramishima/m-backend/internal/adapters/cache/redis"
	"kiramishima/m-backend/internal/adapters/database/postgresql/repository"
	"kiramishima/m-backend/internal/adapters/pubsub/psnats"
	"kiramishima/m-backend/internal/core/services"
	"kiramishima/m-backend/internal/handlers"
	"kiramishima/m-backend/internal/server"

	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func bootstrap(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	server *server.Server,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting API")

				/*go func() {
					if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal("failed to start server")
					}
				}()*/
				_ = server.Run()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return logger.Sync()
			},
		},
	)
}

var Module = fx.Options(
	config.Module,
	config.LoggerModule,
	fx.Provide(func() *chi.Mux {
		var r = chi.NewRouter()
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"*"},
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Use(middleware.Timeout(60 * time.Second))
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		r.Use(httprate.LimitByIP(1000, 1*time.Minute))
		r.Use(middleware.Compress(5))
		return r
	}),
	fx.Provide(func() *render.Render {
		return render.New()
	}),
	fx.Provide(func() *validator.Validate {
		return validator.New(validator.WithRequiredStructEnabled())
	}),
	server.Module,
	repository.DatabaseModule,
	services.Module,
	handlers.Module,
	redis.Module,
	psnats.Module,
	fx.Invoke(bootstrap),
)
