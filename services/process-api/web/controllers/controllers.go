package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type API struct {
}

var _ ServerInterface = (*API)(nil)

func NewAPI() *API {
	return &API{}
}

func (s *API) ListAvailableProcesses(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)

}

// (GET /process/healthz)
func (s *API) GetProcessHealthz(c *gin.Context) {}

// List Jobs
// (GET /process/jobs)
func (s *API) ListJobs(c *gin.Context) {}

// Trigger Process Job
// (POST /process/jobs)
func (s *API) TriggerProcessJob(c *gin.Context) {}

// Cancel Process
// (POST /process/jobs/{jobId}/cancel)
func (s *API) CancelProcess(c *gin.Context, jobId JobId) {}

// Get Process Results
// (GET /process/jobs/{jobId}/results)
func (s *API) GetProcessResults(c *gin.Context, jobId JobId) {}

// Get Process Status
// (GET /process/jobs/{jobId}/status)
func (s *API) GetProcessStatus(c *gin.Context, jobId JobId) {}

// (GET /process/metrics)
func (s *API) GetProcessMetrics(c *gin.Context) {}
func (s *API) RegisterEndpoint(r *gin.Engine, baseURL string, middlewares []MiddlewareFunc) {
	RegisterHandlersWithOptions(r, s, GinServerOptions{
		BaseURL:     baseURL,
		Middlewares: middlewares,
	})
}
