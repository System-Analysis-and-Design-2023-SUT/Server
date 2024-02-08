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

		var port = string(m.Meta)
		response, err := http.Get(fmt.Sprintf("http://%s:%s/_pull?key=%s", m.Addr, port, data.Key))
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}

	return nil
}

func (h *Helper) Write(data models.Data) error {
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var port = string(m.Meta)
		response, err := http.Post(fmt.Sprintf("http://%s:%s/_push?key=%s&value=%s", m.Addr, port, data.Key, data.Value),
			"application/json",
			nil,
		)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}

	return nil
}

func (h *Helper) GetQueue() ([]byte, error) {
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var port = string(m.Meta)
		response, err := http.Get(fmt.Sprintf("http://%s:%s/queue", m.Addr, port))
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	for _, m := range h.list.Members() {
		fmt.Println(m.Name, m.Addr)
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
