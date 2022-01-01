package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/url"
	pb "order-service/pb/generated"
	"os"
	"os/signal"
	"syscall"
	"time"
)


const (
	TCPNetwork = "tcp"
	GrpcAddress = "localhost:7777"
	TCPAddress = "localhost:8888"
	WebSocketAddress = "localhost:9999"
)

func main()  {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	conn, err := net.Dial(TCPNetwork, TCPAddress)
	if err != nil {
		logrus.Errorf("Connect to TCP Server error: %s", err.Error())
		panic(err)
	}

	defer conn.Close()

	data := make([]byte, 4096)
	length, err := conn.Read(data)
	if err != nil {
		logrus.Errorf("Cannot received response from TCP Server: %s", err.Error())
		panic(err)
	}

	response := &pb.GetProductDetailResponse{}
	if err := proto.Unmarshal(data[:length], response); err != nil {
		logrus.Errorf("Cannot parse response from TCP Server: %s", err.Error())
		panic(err)
	}

	log.Printf("Receivce Product Raw data: %s", string(data))
	log.Printf("Receivce Product response: %v", response)


	go func() {
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)

			var addr = flag.String("", WebSocketAddress, "websocket service address")
			productServiceWsUrl := url.URL{Scheme: "ws", Host: *addr, Path: "/products"}
			c, _, err := websocket.DefaultDialer.Dial(productServiceWsUrl.String(), nil)
			if err != nil {
				logrus.Infof("Dial to Product Service WebSocket error: %s", err.Error())
				return
			}
			defer c.Close()

			done := make(chan struct{})
			go func() {
				defer close(done)
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						return
					}
					log.Printf("recv: %s", message)
				}
			}()

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
					if err != nil {
						log.Println("write:", err)
						return
					}
				case <-interrupt:
					log.Println("interrupt")
					err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-done:
					case <-time.After(time.Second):
					}
					return
				}
			}
	}()

	logrus.Error("Service error %s", <-errs)

}
