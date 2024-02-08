package queue

import (
	"bytes"
	"log"
	"net/http"

	repo "github.com/System-Analysis-and-Design-2023-SUT/Server/internal/repository/queue"
	service "github.com/System-Analysis-and-Design-2023-SUT/Server/internal/services/queue"
	logging "github.com/System-Analysis-and-Design-2023-SUT/Server/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var logger *logging.Logger

func init() {
	var err error
	logger, err = logging.NewLogger("server_api_queue", true)
	if err != nil {
		log.Fatal("could not initialize server api queue module logger")
	}
}

type Queue struct {
	repository *repo.Repository
	service    *service.Service
}

func (q *Queue) RegisterRoutes(v1 *gin.RouterGroup) {
	logger.InfoS("Registering queue related endpoints to api server.")

	api := v1.Group("/")

	api.POST("/push", q.pushEndpoint())          // Push into queue.
	api.POST("/_push", q.pushForceEndpoint())    // Force push into queue.
	api.GET("/pull", q.pullEndpoint())           // Gets head of queue.
	api.GET("/_pull", q.pullForceEndpoint())     // Gets head of queue.
	api.GET("/subscribe", q.subscribeEndpoint()) // Subscribe in queue.
	api.GET("/queue", q.copyEndpoint())          // Gets whole of queue.
}

func (q *Queue) pushEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q.service.Push(c, false)
	}
}

func (q *Queue) pushForceEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q.service.Push(c, true)
	}
}

func (q *Queue) pullEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q.service.Pull(c, false)
	}
}

func (q *Queue) pullForceEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q.service.Pull(c, true)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (q *Queue) subscribeEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		addr := conn.RemoteAddr().String()
		if err != nil {
			logger.Error(err.Error())
			return
		}
		defer func() {
			conn.Close()
			err := q.service.Unsubscribe(conn, addr)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			if bytes.Equal(msg, []byte("subscribe\n")) {
				response := q.service.Subscribe(conn, addr)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
					logger.Error(err.Error())
					return
				}

			} else {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("Invalid Message")); err != nil {
					logger.Error(err.Error())
					return
				}
			}
		}

	}
}

func (q *Queue) copyEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q.service.Copy(c)
	}
}

func NewQueueModule(repo *repo.Repository, service *service.Service) (*Queue, error) {
	if repo == nil {
		return nil, ErrNilQueueRepo
	}

	if service == nil {
		return nil, ErrNilQueueService
	}

	return &Queue{
		repository: repo,
		service:    service,
	}, nil
}
