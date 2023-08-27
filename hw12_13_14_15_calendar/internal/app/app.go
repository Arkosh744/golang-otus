package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/handlers"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/handlers/middlewares"
	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/log"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		if err := app.RunHTTPServer(ctx); err != nil {
			log.Fatalf("failed to run http server: %v", err)
		}
	}()

	defer app.StopHTTPServer(ctx)

	<-ctx.Done()

	return nil
}

func (app *App) initDeps(ctx context.Context) error {
	for _, init := range []func(context.Context) error{
		config.Init,
		log.InitLogger,
		app.initServiceProvider,
		app.initHTTPServer,
	} {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initServiceProvider(ctx context.Context) error {
	app.serviceProvider = newServiceProvider(ctx)

	return nil
}

func (app *App) initHTTPServer(_ context.Context) error {
	const timeout = 15

	h := handlers.InitRouter(app.serviceProvider.calendarService)

	app.httpServer = &http.Server{
		Addr:         net.JoinHostPort(config.AppConfig.Host, config.AppConfig.Port),
		Handler:      middlewares.LoggingMiddleware(h),
		ReadTimeout:  timeout * time.Second,
		WriteTimeout: timeout * time.Second,
	}

	return nil
}

func (app *App) RunHTTPServer(_ context.Context) error {
	log.Infof("Start: HTTP server listening on port %s", config.AppConfig.Port)

	if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func (app *App) StopHTTPServer(ctx context.Context) {
	log.Infof("Stop: HTTP server on port %s", config.AppConfig.Port)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		log.Info("Shutdown http server")
		cancel()
	}()

	if err := app.httpServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Info("failed to stop http server: %v", err)
	}
}
