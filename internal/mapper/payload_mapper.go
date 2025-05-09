package mapper

import (
	"github.com/developeerz/restorio-reserving/internal/adapter/postgres/entity"
	"github.com/developeerz/restorio-reserving/internal/pkg/models"
)

func ToPayload(payloadEntity entity.OutboxPayload, reservationTime string, telegramID int) models.PayloadTelegram {
	return models.PayloadTelegram{
		RestaurantName:    payloadEntity.RestaurantName,
		RestaurantAddress: payloadEntity.RestaurantAddress,
		TableNumber:       payloadEntity.TableNumber,
		ReservationTime:   reservationTime,
		TelegramID:        telegramID,
	}
}
