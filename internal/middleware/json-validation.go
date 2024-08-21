package middleware

import (
	"encoding/json"
	"errors"
	"wb/internal/models"
)

func NewOrder(data []byte) (models.Order, error) {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return models.Order{}, err
	}

	err := ValidateOrder(&order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func ValidateOrder(order *models.Order) error {

	var field interface{} = order.OrderUid

	if order.OrderUid == "" {
		return errors.New("invalid order uid")
	} else if _, ok := field.(string); !ok {
		return errors.New("invalid type of order uid")
	}

	field = order.TrackNumber
	if order.TrackNumber == "" {
		return errors.New("invalid order track number")
	} else if _, ok := field.(string); !ok {
		return errors.New("invalid type of order track number")
	}

	field = order.Payment.Transaction
	if order.Payment.Transaction == "" {
		return errors.New("invalid transaction")
	} else if _, ok := field.(string); !ok {
		return errors.New("invalid type of order transaction")
	}

	field = order.Payment.Amount
	if order.Payment.Amount == 0 {
		return errors.New("invalid payment amount")
	} else if _, ok := field.(int); !ok {
		return errors.New("invalid type of order amount")
	}

	return nil

}
