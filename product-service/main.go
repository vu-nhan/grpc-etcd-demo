package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"product-service/handlers"
	pb "product-service/pb/generated"
	"syscall"
)

const (
	TCPNetwork = "tcp"
	GrpcAddress = "localhost:7777"
	TCPAddress = "localhost:8888"
	WebSocketAddress = "localhost:9999"
)

func main()  {
	grpcListener, err := net.Listen(TCPNetwork, GrpcAddress)
	if err != nil {
		logrus.Errorf("Listener GRPC in Address has error: %s", err.Error())
		panic(err)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcProductHandler := handlers.NewGrpcProductHandler()
	go func() {
		grpcServer := grpc.NewServer()
		pb.RegisterProductServiceServer(grpcServer, grpcProductHandler)

		logrus.Info("Start GRPC Transport Server")
		if err := grpcServer.Serve(grpcListener); err != nil {
			logrus.Errorf("Start in GRPC Server error: %s", err.Error())
			os.Exit(1)
		}
	}()

	tcpListener, err := net.Listen(TCPNetwork, TCPAddress)
	if err != nil {
		logrus.Errorf("Listener TCP in Address has error: %s", err.Error())
		panic(err)
	}

	tcpProductHandler := handlers.NewTCPProductHandler()
	go func() {
		for {
			if conn, err := tcpListener.Accept(); err == nil {
				tcpProduct, _ := tcpProductHandler.GetProductDetail(context.Background(), nil)
				data, err := proto.Marshal(tcpProduct)
				if err != nil {
					logrus.Errorf("Marshal TCP Product Data to ProtoBuf error: %s", err.Error())
				}

				length, err := conn.Write(data)
				if err != nil {
					logrus.Errorf("Return TCP Product Data error: %s", err.Error())
				}
				logrus.Info("Length: %d", length)
				logrus.Info("Return Product successfully")
			} else {
				logrus.Errorf("Listen TCP error: %s", err.Error())
			}
		}
	}()

	go func() {
		var upgrader = websocket.Upgrader{}
		var addr = flag.String("", WebSocketAddress, "WebSocket address")
		http.HandleFunc("/products", func(writer http.ResponseWriter, request *http.Request) {
			connection, err := upgrader.Upgrade(writer, request, nil)
			if err != nil {
				log.Print("upgrade:", err)
				return
			}
			defer connection.Close()
			for {
				messageType, message, err := connection.ReadMessage()
				if err != nil {
					logrus.Errorf("Read Message error: %s", err.Error())
					break
				}
				logrus.Infof("Read Message successfully: %s", string(message))

				request := &pb.GetProductDetailRequest{}
				if err := proto.Unmarshal(message, request); err != nil {
					logrus.Errorf("Cannot parse response from TCP Server: %s", err.Error())
					panic(err)
				}

				webSocketProductHandler := handlers.NewWebSocketProductHandler()
				response := webSocketProductHandler.GetProductDetail(context.Background(), "")
				log.Printf("Recevie Message from order-service: %s", message)
				err = connection.WriteMessage(messageType, getByte(response))
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		})
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	logrus.Error("Service error %s", <-errs)
}

func getByte(data *pb.GetProductDetailResponse) []byte {
	response, _ := json.Marshal(data)
	return response
}