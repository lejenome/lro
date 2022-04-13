package server

import (
	"github.com/lejenome/lro/services/process-api/web/controllers"
	c "github.com/lejenome/lro/services/process-api/web/controllers"
	"github.com/lejenome/lro/services/process-api/web/middlewares"
)

func (s *Server) setupRoutes() {
	AuthMiddleware := middlewares.NewJWTAuthenticationMiddleware(&s.auth.JWT).Middleware()
	{
		api := c.NewAPI()
		api.RegisterEndpoint(s.routes, "/api", []controllers.MiddlewareFunc{
			controllers.MiddlewareFunc(AuthMiddleware),
		})
	}
}
