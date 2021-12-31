package handlers

import (
	"context"
	pb "order-service/pb/generated"
)

type TcpOrderHandler struct {

}

func NewTcpOrderHandler() *TcpOrderHandler {
	return &TcpOrderHandler{}
}

func (h *TcpOrderHandler) GetProductDetail(ctx context.Context, productId string) *pb.GetProductDetailResponse {
	request := &pb.GetProductDetailRequest{ProductId: productId}


}
