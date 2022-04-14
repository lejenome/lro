package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-api/web/controllers"

	"github.com/gin-gonic/gin"
)

// Server is the main server component, it holds context objects and exec the
// server
type Server struct {
	routes *gin.Engine
	auth   config.AuthConfig
	ctx    context.Context
	srv    *http.Server
}

// New create new instance of the Server
func New(ctx context.Context, env string, auth config.AuthConfig) Server {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	return Server{
		routes: gin.New(),
		auth:   auth,
		ctx:    ctx,
		srv:    nil,
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

func (s *Server) ListenAndServe(addr string) error {
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.routes,
	}
	go func() {
		<-s.ctx.Done()
		s.Shutdown()
	}()
	err := s.srv.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	if s.srv == nil {
		return
	}
	// the server has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	s.srv = nil
}
