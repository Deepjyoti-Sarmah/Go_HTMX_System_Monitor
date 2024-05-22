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
	"github.com/rs/cors"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscriberMutex         sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() (*server, http.Handler) {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribHandler)

	 corsConfig := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:8080", "https://webtop-pw7d.onrender.com"},
        AllowCredentials: true,
        Debug:            true,
    })

    handler := corsConfig.Handler(&s.mux)
	// handler := cors.Default().Handler(&s.mux)
	return s, handler
}

func (s *server) subscribHandler(writer http.ResponseWriter, req *http.Request) {
	err := s.subscribe(req.Context(), writer, req)
	if err != nil {
		fmt.Println("err",err)
		return
	}
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)
	defer s.removeSubscriber(subscriber)

	c, err := websocket.Accept(writer, req, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancle := context.WithTimeout(ctx, time.Second)
			defer cancle()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *server) removeSubscriber(subscriber *subscriber) {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()
	delete(s.subscribers, subscriber)
	fmt.Println("Removed subscriber", subscriber)
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
	srv, handler := NewServer()
	go func(s *server) {
		for {
			systemData, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err)
				continue
			}

			diskData, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err)
				continue
			}

			cpuData, err := hardware.GetCpuSection()
			if err != nil {
				fmt.Println(err)
				continue
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			msg := []byte(`
				<div hx-swap-oob="innerHTML:#update-timestamp">
					<p><i style="color: green" class="fa fa-circle"></i> ` + timeStamp + `</p>
				</div>
				<div hx-swap-oob="innerHTML:#system-data">` + systemData + `</div>
				<div hx-swap-oob="innerHTML:#cpu-data">` + cpuData + `</div>
				<div hx-swap-oob="innerHTML:#disk-data">` + diskData + `</div>
			`)
			s.publishMsg(msg)
			time.Sleep(3 * time.Second)
		}
	}(srv)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
