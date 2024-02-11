package helper

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/settings"
	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	"github.com/hashicorp/memberlist"
)

type Helper struct {
	list *memberlist.Memberlist
	st   *settings.Settings
}

func (h *Helper) Read(data models.Data) error {
	ips, er := net.LookupIP(h.st.Replica.Hostname[0])
	if er != nil {
		fmt.Println("Error:", er)
	}

	ip, ipnet, err := net.ParseCIDR(h.st.Replica.Subnet)
	if err != nil {
		fmt.Println("Error parsing subnet:", err)
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if ip.Equal(ips[0]) {
			continue
		}

		response, err := http.Get(fmt.Sprintf("http://%s:8080/_pull?key=%s", ip, data.Key))
		if err != nil {
			continue
		}

		defer response.Body.Close()
	}
	return nil
}

func (h *Helper) Write(data models.Data) error {
	ips, er := net.LookupIP(h.st.Replica.Hostname[0])
	if er != nil {
		fmt.Println("Error:", er)
	}

	ip, ipnet, err := net.ParseCIDR(h.st.Replica.Subnet)
	if err != nil {
		fmt.Println("Error parsing subnet:", err)
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if ip.Equal(ips[0]) {
			continue
		}

		first := false
		for _, m := range h.list.Members() {
			if string(m.Meta) == ip.String() && len(m.Name) > 6 && m.Name[:len(m.Name)-4] == "sad-server-1" {
				first = true
			}
		}
		if first {
			continue
		}

		response, _ := http.Post(fmt.Sprintf("http://%s:8080/_push?key=%s&value=%s", ip, data.Key, data.Value),
			"application/json",
			nil,
		)

		defer response.Body.Close()
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if ip.Equal(ips[0]) {
			continue
		}

		first := false
		for _, m := range h.list.Members() {
			if string(m.Meta) == ip.String() && len(m.Name) > 6 && m.Name[:len(m.Name)-4] == "sad-server-1" {
				first = true
			}
		}
		if !first {
			continue
		}

		response, _ := http.Post(fmt.Sprintf("http://%s:8080/_push?key=%s&value=%s", ip, data.Key, data.Value),
			"application/json",
			nil,
		)

		defer response.Body.Close()
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

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
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

func NewHelper(list *memberlist.Memberlist, st *settings.Settings) (*Helper, error) {
	if list == nil {
		return nil, ErrNilMemberlist
	}

	return &Helper{
		list: list,
		st:   st,
	}, nil
}
