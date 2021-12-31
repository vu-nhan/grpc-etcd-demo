package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	pb "order-service/pb/generated"
)

func main()  {
	conn, err := net.Dial("tcp", "localhost:8888")
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
}
