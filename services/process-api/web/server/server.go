package server

import (
	"net/http"

	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-api/web/controllers"

	"github.com/gin-gonic/gin"
)

// Server is the main server component, it holds context objects and exec the
// server
type Server struct {
	routes *gin.Engine
	auth   config.AuthConfig
}

// New create new instance of the Server
func New(env string, auth config.AuthConfig) Server {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	return Server{
		gin.New(),
		auth,
	}
}

func (s *Server) registerEndpoint(group *gin.RouterGroup, path string, controller controllers.Controller) {
	var r *gin.RouterGroup
	if len(path) > 0 {
		r = group.Group(path)
	} else {
		r = group
	}
	controller.RegisterEndpoint(r)
}

func (s *Server) Setup() {
	s.setupMiddlewares()
	s.setupRoutes()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.routes.ServeHTTP(w, req)
}

func (s *Server) ListenAndServe(addr string) {
	s.routes.Run(addr)
}
