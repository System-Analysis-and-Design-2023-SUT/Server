package queue

import (
	"fmt"

	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/helper"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/settings"
	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Repository struct {
	st         *settings.Settings
	helper     *helper.Helper
	queue      *models.Queue
	subscriber *models.Subscriber
}

func NewRepository(st *settings.Settings, helper *helper.Helper, q *models.Queue, s *models.Subscriber) (*Repository, error) {
	if st == nil {
		return nil, errors.New("st should not be nil")
	}
	if helper == nil {
		return nil, errors.New("helper should not be nil")
	}
	if q == nil {
		return nil, errors.New("queue should not be nil")
	}

	d, err := helper.GetQueue()
	if err != nil {
		fmt.Println(err)
	} else {
		err := q.BulkPush(d)
		return &Repository{}, err
	}

	return &Repository{
		st:         st,
		helper:     helper,
		queue:      q,
		subscriber: s,
	}, nil
}

// Push will save data into queue
func (r *Repository) Push(data models.Data, force bool) (models.Data, error) {
	if len(r.subscriber.Member) > 0 {
		err := r.helper.Read(data)
		if err != nil {
			return models.Data{}, err
		}

		err = r.subscriber.Send(data)
		return data, err
	}
	err := r.queue.Push(data)
	if err != nil {
		return models.Data{}, err
	}
	if !force {
		err = r.helper.Write(data)
		if err != nil {
			return models.Data{}, err
		}
	}
	return data, nil
}

// Pull return head of queue
func (r *Repository) Pull(key string) (models.Data, error) {
	if key != "" {
		return models.Data{}, r.queue.Delete(key)
	}
	d, err := r.queue.Pull()
	if err != nil {
		return models.Data{}, err
	}

	err = r.helper.Read(d)
	if err != nil {
		return models.Data{}, err
	}
	return d, nil
}

func (r *Repository) Subscribe(c *websocket.Conn, addr string) string {
	resp, err := r.subscriber.Subscribe(c, addr)
	if err != nil {
		return err.Error()
	}
	return resp
}

func (r *Repository) Unsubscribe(c *websocket.Conn, addr string) error {
	return r.subscriber.Unsubscribe(addr)
}

// Copy return whole of queue
func (r *Repository) Copy() *models.Queue {
	return r.queue
}
