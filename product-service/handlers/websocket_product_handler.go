package handlers

import (
	"context"
	"github.com/google/uuid"
	pb "product-service/pb/generated"
	"time"
)

type WebSocketProductHandler struct {
}

func (h *WebSocketProductHandler) GetProductDetail(ctx context.Context, productId string) (*pb.GetProductDetailResponse) {
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
	}
}

func NewWebSocketProductHandler() *WebSocketProductHandler {
	return &WebSocketProductHandler{}
}

