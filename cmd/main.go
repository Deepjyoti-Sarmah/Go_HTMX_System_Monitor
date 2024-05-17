package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/deepjyoti-sarmah/go_htmx_hardware_monitor/internal/hardware"
	"nhooyr.io/websocket"
)

type server struct {
	subscriberMessageBuffer int 
	mux http.ServeMux
	subscriberMutex sync.Mutex
	subscribers map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers: make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribHandler)
	return s
}

func (s *server) subscribHandler(writer http.ResponseWriter, req *http.Request)  {
	err := s.subscribe(req.Context(), writer, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)

	c, err := websocket.Accept(writer, req, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <- subscriber.msgs:
			ctx, cancle := context.WithTimeout(ctx, time.Second)
			defer cancle()
			err := c.Write(ctx, websocket.MessageText, msg) 
			if err != nil {
				return err
			}
		case <- ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscriberMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscriberMutex.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

func (cs *server) publishMsg(msg []byte) {
	cs.subscriberMutex.Lock()
	defer cs.subscriberMutex.Unlock()

	for s := range cs.subscribers {
		s.msgs <- msg
	}
}

func main() {
	fmt.Println("Starting system monitor...")
	srv := NewServer()
	go func (s *server)  {
		for {
			systemSection, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err)
			}

			diskSection, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err)
			}

			cpuSection, err := hardware.GetCpuSection()
			if err != nil {
				fmt.Println(err)
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			html := `
			<div hx-swap-oob="innerHTML:#update-timestamp">
				`+timeStamp+`
			</div>
			<div hx-swap-oob="innerHTML:#system-data">
				`+systemSection+`
			</div>
			<div hx-swap-oob="innerHTML:#disk-data">
				`+diskSection+`
			</div>
			<div hx-swap-oob="innerHTML:#cpu-data">
				`+cpuSection+`
			</div>
			`
			s.publishMsg([]byte(html))

			time.Sleep(3 * time.Second)
		}
	}(srv)

	err := http.ListenAndServe(":8080", &srv.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
