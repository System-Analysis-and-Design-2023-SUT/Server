package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Health struct {
	Environment string
}

// ReadinessRequest Input Model
type ReadinessRequest struct {
}

// ReadinessResponse represents body of Readiness API.
type ReadinessResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// LivenessRequest Input Model
type LivenessRequest struct {
}

// LivenessResponse represents body of Liveness API.
type LivenessResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (h *Health) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/-/ready", h.isReady)
	group.GET("/-/live", h.isLive)
}

func (h *Health) isReady(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"reason": "",
	})
}

func (h *Health) isLive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"reason": "",
	})
}

func NewHealth(env string) (*Health, error) {
	return &Health{
		Environment: env,
	}, nil
}
