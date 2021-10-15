// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"

	"github.com/evrone/go-clean-template/config"
	amqprpc "github.com/evrone/go-clean-template/internal/controller/amqp_rpc"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/repo"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/mysql"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html.
// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := mysql.New(cfg.Mysql.URL, mysql.MaxPoolSize(cfg.Mysql.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	translationUseCase := usecase.New(
		repo.New(pg),
		webapi.New(),
	)

	// RabbitMQ RPC Server, for remote call
	rmqRouter := amqprpc.NewRouter(translationUseCase)
	rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, "fanout", rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, translationUseCase)

	// service info for etcd or no
	_ = httpserver.NewNoEtcd(handler, net.JoinHostPort("", cfg.HTTP.Port))

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

	l.Info("app - Run - exit ! ")
}
