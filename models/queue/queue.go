package models

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Member map[string]*websocket.Conn
	List   []string
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Queue struct {
	KeySet map[string]struct{} `json:"keySet"`
	List   []Data              `json:"list"`
}

func (q *Queue) Push(data Data) error {
	fmt.Println("TTTTTTTTTT")
	fmt.Println(data)
	fmt.Println(q)
	fmt.Println(q.KeySet)
	fmt.Println(q.KeySet[data.Key])

	if _, ok := q.KeySet[data.Key]; ok {
		return ErrKeyExist
	}
	q.KeySet[data.Key] = struct{}{}

	q.List = append(q.List, data)
	return nil
}

func (q *Queue) Pull() (Data, error) {
	if len(q.List) == 0 {
		return Data{}, ErrEmptyList
	}

	delete(q.KeySet, q.List[0].Key)
	result := q.List[0]
	q.List = q.List[1:]
	return result, nil
}

func (q *Queue) Delete(key string) error {
	if _, ok := q.KeySet[key]; !ok {
		return ErrKeyNotFound
	}

	delete(q.KeySet, key)
	for i, l := range q.List {
		if l.Key == key {
			if i+1 == len(q.List) {
				q.List = q.List[:i]
			} else {
				q.List = append(q.List[:i], q.List[i+1:]...)
			}
			return nil
		}
	}
	return ErrObjectNotFound
}

func (q *Queue) BulkPush(data []byte) error {
	var tmp Queue

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return ErrParseData
	}

	for _, l := range tmp.List {
		err := q.Push(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewQueue() *Queue {
	return &Queue{
		KeySet: make(map[string]struct{}),
		List:   make([]Data, 0),
	}
}

func (s *Subscriber) Subscribe(c *websocket.Conn, addr string) (string, error) {
	if _, ok := s.Member[addr]; ok {
		return "", ErrSubscriberExist
	}
	s.Member[addr] = c
	s.List = append(s.List, addr)
	return "You subscribe successfully", nil
}

func (s *Subscriber) Unsubscribe(addr string) error {
	delete(s.Member, addr)
	for i, l := range s.List {
		if l == addr {
			if i+1 == len(s.List) {
				s.List = s.List[:i]
			} else {
				s.List = append(s.List[:i], s.List[i+1:]...)
			}
			return nil
		}
	}
	return nil
}

func (s *Subscriber) Send(data Data) error {
	index := rand.Intn(len(s.Member))
	conn := s.Member[s.List[index]]
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(body))
	return err
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Member: make(map[string]*websocket.Conn),
	}
}
