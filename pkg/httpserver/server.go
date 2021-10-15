// Package httpserver implements HTTP server.
package httpserver

import (
	"common/monitor"
	"common/registry"
	"common/service-wrapper"
	"fmt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	_defaultAddr            = ":80"
)

var (
	// for etcd .
	Service = &registry.Service{
		Name: "go.micro.api.media-proxy",
		Metadata: map[string]string{
			"serverDescription": "audio recording proxy service", // server desc.
		},
		Nodes: []*registry.Node{
			{
				Id:      "go.micro.api.media-proxy-",
				Address: "localhost",
				Port:    8400,
				Metadata: map[string]string{
					"serverTag":           "media-proxy", // server division.
					monitor.ServiceStatus: monitor.DeleteState,
				},
			},
		},
		Version: "2",
	}
)

// Server -.
type Server struct {
	service        service_wrapper.Service
}

// New -.
// use self HttpServer.
func NewNoEtcd(handler *gin.Engine, l logger.Interface, address string) *Server {
	service := service_wrapper.NewService(service_wrapper.Address(address),
		service_wrapper.Engine(handler))
	// for test no etcd.
	// service_wrapper.ServiceInfo(Service))
	// service_wrapper.RegisterInterval(monitor.HeartBeatCheck)) , service_wrapper.Registry(registry.DefaultRegistry)) no etcd
	if err := service.Run(); err != nil {
		_ = fmt.Errorf("service StartFail:%v", err)
		l.Error("service StartFail:%v", err)
	}
	return &Server{ service }
}

// Shutdown -.
func (s *Server) Shutdown() error {
	return  s.service.Shutdown()
}
