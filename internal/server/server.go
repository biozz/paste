package server

import (
	"net/http"
	"time"

	"github.com/biozz/paste/internal/config"
	"github.com/biozz/paste/internal/storage_types"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Web struct {
	conf            config.Server
	storage         storage_types.Storage
	echo            *echo.Echo
	templatesLoader *TemplatesLoader
}

func New(conf config.Server, storage storage_types.Storage) Web {
	e := echo.New()
	dev := conf.Env == "dev"
	templatesLoader := NewTemplatesLoader(dev)
	e.Renderer = NewRenderer(dev, templatesLoader)
	return Web{
		conf:            conf,
		storage:         storage,
		echo:            e,
		templatesLoader: templatesLoader,
	}
}

func (h *Web) Start() {
	h.echo.Use(middleware.Logger())
	prom := prometheus.NewPrometheus("wow", nil)
	prom.Use(h.echo)

	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 1, Burst: 5, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	rateLimiter := middleware.RateLimiterWithConfig(rateLimiterConfig)

	h.echo.GET("/", h.Index)
	h.echo.POST("/", h.New, rateLimiter)
	h.echo.GET("/:slug", h.View)
	h.echo.GET("/url/:slug", h.Redirect)
	h.echo.GET(
		"/static/*",
		echo.WrapHandler(http.StripPrefix("/static", h.Static())),
	)

	h.echo.Logger.Fatal(h.echo.Start(h.conf.Bind))
}
