package handlers

import (
	"balance-service/dto/proto"
	"balance-service/services/pg"
	log "github.com/sirupsen/logrus"
	"strings"
)

type UpdateOrderEventHandler struct {
	pgCl *pg.Client
}

func NewUpdateOrderEventHandler(pgCl *pg.Client) *UpdateOrderEventHandler {
	return &UpdateOrderEventHandler{
		pgCl: pgCl,
	}
}

func (h *UpdateOrderEventHandler) HandleMessage(body []byte) {
	event, err := unmarshalRequest[*proto.OrderUpdateEvent](body)
	if err != nil {
		log.Errorf("Failed deserialize OrderUpdateEvent: %v", err)
		return
	}
	h.processing(event)
}

type User struct {
	UserId             string
	Balance            float64
	LockedBalance      float64
	CurrencyName       string
	LockedCurrencyName string
}

func (h *UpdateOrderEventHandler) processing(event *proto.OrderUpdateEvent) {
	owner := &User{UserId: event.Order.UserId}
	matchedUser := &User{UserId: event.MatchedUser.UserId}
	base, quote, _ := strings.Cut(event.Order.Pair, "/")

	if event.Order.Status == proto.OrderStatus_ORDER_STATUS_NEW {
		switch event.Order.Direction {
		case proto.Direction_ORDER_DIRECTION_BUY:
			owner.LockedCurrencyName = quote
			owner.LockedBalance = event.Order.InitPrice * event.Order.InitVolume
			break
		case proto.Direction_ORDER_DIRECTION_SELL:
			owner.LockedCurrencyName = base
			owner.LockedBalance = event.Order.InitVolume
			break
		default:
			log.Errorf("unsupported direction value")
			return
		}
	}

	if event.Order.Status == proto.OrderStatus_ORDER_STATUS_MATCHED {
		switch event.Order.Direction {
		case proto.Direction_ORDER_DIRECTION_BUY:
			// Buyer:
			owner.LockedCurrencyName = quote
			owner.LockedBalance = event.MatchedUser.Volume

			owner.CurrencyName = base
			owner.Balance = event.MatchedUser.Volume

			// Seller:
			matchedUser.LockedCurrencyName = base
			matchedUser.LockedBalance = event.MatchedUser.Volume

			matchedUser.CurrencyName = quote // got money
			matchedUser.Balance = event.MatchedUser.Volume * event.MatchedUser.Price
		case proto.Direction_ORDER_DIRECTION_SELL:
			// Seller:
			owner.LockedCurrencyName = base
			owner.LockedBalance = event.MatchedUser.Volume

			owner.CurrencyName = quote
			owner.Balance = event.MatchedUser.Volume * event.MatchedUser.Price

			// Buyer:
			matchedUser.LockedCurrencyName = quote
			matchedUser.LockedBalance = event.MatchedUser.Volume * event.MatchedUser.Price

			matchedUser.CurrencyName = base
			matchedUser.Balance = event.MatchedUser.Volume
		}
	}

	if err := h.pgCl.UpdateBalancesTx(owner, matchedUser); err != nil {
		log.Errorf("failed to update balances: %v", err)
	}
}
