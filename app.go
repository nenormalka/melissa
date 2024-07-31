package melissa

import (
	"context"
	"fmt"
	log "github.com/nenormalka/melissa/logger"
	"time"

	"github.com/nenormalka/melissa/types"
)

const (
	shutdownTimeout = 10 * time.Second
)

type (
	App struct {
		servers  *types.ServerPool
		services *types.ServicePool
	}
)

func NewApp(
	servers *types.ServerPool,
	services *types.ServicePool,
) *App {
	return &App{
		servers:  servers,
		services: services,
	}
}

func (c *App) Run(ctx context.Context) error {
	logger := log.NewLogger()

	logger.Info("Services start")
	if err := c.services.Start(ctx); err != nil {
		return fmt.Errorf("services.Start: %w", err)
	}

	logger.Info("Servers start")
	if err := c.servers.Start(ctx); err != nil {
		return fmt.Errorf("servers.Start: %w", err)
	}

	logger.Info("Application is ready 🐣")

	<-ctx.Done()

	sdCtx, sdCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer sdCancel()

	logger.Info("Stopping servers...")
	c.servers.Stop(sdCtx)

	logger.Info("Stopping services...")
	c.services.Stop(sdCtx)

	logger.Info("Gracefully stopped, bye bye 👋")

	return nil
}
