package dto

import "github.com/google/uuid"

type OrderEstimateResponseDto struct {
	TotalPrice                     int       `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int       `json:"estimatedDeliveryTimeInMinutes"`
	OrderId                        uuid.UUID `json:"calculatedEstimateId"`
}