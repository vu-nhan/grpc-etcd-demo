package handlers

import (
	"context"
	pb "product-service/pb/generated"
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
			Id:          "1234567",
			Code:        "Product-code",
			Name:        "Product-name",
			Description: "Product-description",
			Status:      "ACTIVE",
			CreatedDate: "2021-01-01",
			UpdatedDate: "2021-01-01",
		},
	}, nil
}

func NewTCPProductHandler() *TCPProductHandler {
	return &TCPProductHandler{}
}
