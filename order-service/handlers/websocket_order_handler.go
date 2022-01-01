package handlers

import (
	"context"
	pb "order-service/pb/generated"
)

type WebSocketOrderHandler struct {

}

func NewWebSocketOrderHandler() *WebSocketOrderHandler {
	return &WebSocketOrderHandler{}
}

func (h *WebSocketOrderHandler) GetProductDetail(ctx context.Context, productId string) *pb.GetProductDetailResponse {
	return nil
}
