package middlewares

import (
	"net/http"

	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-api/providers/auth"
	"github.com/lejenome/lro/services/process-api/utils"

	"github.com/gin-gonic/gin"
)

type JWTAuthenticationMiddleware struct {
	AuthMaker auth.Maker
}

func NewJWTAuthenticationMiddleware(config *config.JWTConfig) *JWTAuthenticationMiddleware {
	return &JWTAuthenticationMiddleware{
		auth.NewJWTMaker(config),
	}
}
func (m *JWTAuthenticationMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.AuthMaker.ExtractToken(c.Request)
		var payload *auth.AuthPayload
		if err == nil {
			payload, err = m.AuthMaker.VerifyToken(token)
		}

		if err != nil {
			c.Abort()
			utils.ErrorResponse(c, http.StatusUnauthorized, "Request Unauthorized")
			return
		}

		c.Set("context", payload.Context)
	}
}
