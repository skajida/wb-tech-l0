package getorderinfo

import (
	"encoding/json"
	"net/http"
	"wb-tech-l0/internal/model"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type OrderInfoHandler struct {
	logger  *zap.Logger
	service cacher
}

func NewHandler(logger *zap.Logger, service cacher) *OrderInfoHandler {
	return &OrderInfoHandler{logger: logger, service: service}
}

func (oih OrderInfoHandler) GetOrderInfoHandle(w http.ResponseWriter, r *http.Request) {
	orderUid, exists := mux.Vars(r)["order_uid"]
	if !exists {
		oih.logger.Error("Handle error: incorrect order_uid request")
		http.Error(w, "order_uid is invalid", http.StatusBadRequest)
		return
	}

	order, err := oih.service.GetOrderInfo(r.Context(), model.OrderId(orderUid))
	switch err {
	case nil:
	case model.ErrOrderBadParam:
		oih.logger.Error("Get order bad parameter", zap.Error(err))
		http.Error(w, "No such order", http.StatusBadRequest)
		return
	default:
		oih.logger.Error("Internal error", zap.Error(err))
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
