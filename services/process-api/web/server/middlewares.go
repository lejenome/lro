package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupMiddlewares() {
	s.routes.Use(gin.Logger())
	s.routes.Use(gin.Recovery())
	s.routes.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"https://foo.com"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Authorization", "Content-Type"},
		// Headers allowed to be accessed by JS from the browser
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowWebSockets:        true,
		AllowBrowserExtensions: false,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))
}
