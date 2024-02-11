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
	fmt.Println("PULLLL")
	fmt.Println(h.list.Members())
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}

		var address = string(m.Meta)
		response, err := http.Get(fmt.Sprintf("http://%s:8080/_pull?key=%s", address, data.Key))
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer response.Body.Close()
	}
	return nil
}

func (h *Helper) Write(data models.Data) error {
	fmt.Println("PUSHHHH")
	fmt.Println(h.list.Members())
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}
		if len(m.Name) > 6 && m.Name[:len(m.Name)-4] == "sad-server-1" {
			continue
		}

		var address = string(m.Meta)
		response, err := http.Post(fmt.Sprintf("http://%s:8080/_push?key=%s&value=%s", address, data.Key, data.Value),
			"application/json",
			nil,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer response.Body.Close()
	}
	for _, m := range h.list.Members() {
		if m == h.list.LocalNode() {
			continue
		}
		if len(m.Name) > 6 && m.Name[:len(m.Name)-4] == "sad-server-1" {
			var address = string(m.Meta)
			response, err := http.Post(fmt.Sprintf("http://%s:8080/_push?key=%s&value=%s", address, data.Key, data.Value),
				"application/json",
				nil,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}

			defer response.Body.Close()
		}
	}
	return nil
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

func (h *Helper) GetFirst() string {
	for _, m := range h.list.Members() {
		if len(m.Name) > 6 && m.Name[:len(m.Name)-4] == "sad-server-1" {
			return string(m.Meta)
		}
	}
	fmt.Println("Can not get first node")
	return ""
}

func NewHelper(list *memberlist.Memberlist) (*Helper, error) {
	if list == nil {
		return nil, ErrNilMemberlist
	}

	return &Helper{
		list: list,
	}, nil
}
