package models

import "time"

type Payment struct {
	OrderUid     string    `json:"order_uid"`
	Transaction  string    `json:"transaction"`
	RequestId    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDt    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}
