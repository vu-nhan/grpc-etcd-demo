package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"product-service/handlers"
	pb "product-service/pb/generated"
	"syscall"
)

func main()  {
	grpcListener, err := net.Listen("tcp", "localhost:8081")
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

	tcpListener, err := net.Listen("tcp", "localhost:8888")
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

	logrus.Error("Service error %s", <-errs)

}