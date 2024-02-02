package queue

import (
	"errors"
	"net/http"

	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/repository/queue"
	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service struct {
	repo *queue.Repository
}

func NewService(repo *queue.Repository) (*Service, error) {
	if repo == nil {
		return nil, errors.New("queue repository should not be nil")
	}
	return &Service{
		repo: repo,
	}, nil
}

func (s *Service) Push(c *gin.Context, force bool) {
	key := c.Query("key")
	value := c.Query("value")

	resp, err := s.repo.Push(
		models.Data{
			Key:   key,
			Value: value,
		},
		force,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *Service) Pull(c *gin.Context, force bool) {
	key := c.Query("key")
	if !force {
		key = ""
	}

	resp, err := s.repo.Pull(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *Service) Subscribe(c *websocket.Conn, addr string) string {
	return s.repo.Subscribe(c, addr)
}

func (s *Service) Unsubscribe(c *websocket.Conn, addr string) error {
	return s.repo.Unsubscribe(c, addr)
}

func (s *Service) Copy(c *gin.Context) {
	c.JSON(http.StatusOK, s.repo.Copy())
}
