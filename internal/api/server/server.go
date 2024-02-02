package server

import (
	"net/http"

	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api/health"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api/queue"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	environment string
	engine      *gin.Engine
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}

func NewServer(queue *queue.Queue, healthMod *health.Health, settings *settings.Settings) (*Server, error) {
	if healthMod == nil {
		return nil, ErrNilHealthModule
	}

	if queue == nil {
		return nil, ErrNilQueueModule
	}

	gin.SetMode(settings.Global.Environment) //todo
	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials:       true,
		ExposeHeaders:          []string{"Content-Length", "Content-Type"},
		AllowBrowserExtensions: true,
	}))

	v1 := engine.Group("/")
	healthMod.RegisterRoutes(v1)
	queue.RegisterRoutes(v1)

	return &Server{
		environment: settings.Global.Environment,
		engine:      engine,
	}, nil
}
