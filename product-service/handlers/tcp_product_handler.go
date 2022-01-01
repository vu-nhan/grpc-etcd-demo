package handlers

import (
	"context"
	"github.com/google/uuid"
	pb "product-service/pb/generated"
	"time"
)

type TCPProductHandler struct {
	pb.UnimplementedProductServiceServer
}

func (h *TCPProductHandler) GetProductDetail(ctx context.Context, request *pb.GetProductDetailRequest) (*pb.GetProductDetailResponse, error) {
	return &pb.GetProductDetailResponse{
		Meta: &pb.Meta{
			Code:    "1",
			Message: "Successfully",
		},
		Data: &pb.Product{
			Id:          uuid.New().String(),
			Code:        "Product-code",
			Name:        "Product-name",
			Description: "Product-description",
			Status:      "ACTIVE",
			CreatedDate: time.Now().String(),
			UpdatedDate: time.Now().String(),
		},
	}, nil
}

func NewTCPProductHandler() *TCPProductHandler {
	return &TCPProductHandler{}
}
