package helper

import (
	"fmt"
	"io"
	"net/http"

	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	"github.com/hashicorp/memberlist"
)

type Helper struct {
	list *memberlist.Memberlist
}

func (h *Helper) Read(data models.Data) error {
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var address = string(m.Meta)
		response, err := http.Get(fmt.Sprintf("http://%s:8080/_pull?key=%s", address, data.Key))
		if err != nil {
			continue
		}

		defer response.Body.Close()
		return nil
	}

	return ErrNodesAreNotReachable
}

func (h *Helper) Write(data models.Data) error {
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var address = string(m.Meta)
		response, err := http.Post(fmt.Sprintf("http://%s:8080/_push?key=%s&value=%s", address, data.Key, data.Value),
			"application/json",
			nil,
		)
		if err != nil {
			continue
		}

		defer response.Body.Close()
		return nil
	}

	return ErrNodesAreNotReachable
}

func (h *Helper) GetQueue() ([]byte, error) {
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var address = string(m.Meta)
		response, err := http.Get(fmt.Sprintf("http://%s:8080/queue", address))
		if err != nil {
			continue
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			continue
		}

		return body, nil
	}

	return nil, ErrQueueNotFound
}

func NewHelper(list *memberlist.Memberlist) (*Helper, error) {
	if list == nil {
		return nil, ErrNilMemberlist
	}

	return &Helper{
		list: list,
	}, nil
}
