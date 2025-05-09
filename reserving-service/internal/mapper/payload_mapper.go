package mapper

import (
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/outbox"
	"github.com/developeerz/restorio-reserving/reserving-service/pkg/models"
)

func ToPayload(payloadEntity outbox.Payload, reservationTime string, telegramID int) models.PayloadTelegram {
	return models.PayloadTelegram{
		RestaurantName:    payloadEntity.RestaurantName,
		RestaurantAddress: payloadEntity.RestaurantAddress,
		TableNumber:       payloadEntity.TableNumber,
		ReservationTime:   reservationTime,
		TelegramID:        telegramID,
	}
}
