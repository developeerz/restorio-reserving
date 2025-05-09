package entity

import "encoding/json"

// OutboxPayload — данные, которые мы читаем в getTablePayload
type OutboxPayload struct {
	RestaurantName    string `db:"name"`
	RestaurantAddress string `db:"address"`
	TableNumber       string `db:"table_number"`
}

// PayloadToJSON преобразует OutboxPayload в JSON ([]byte)
func PayloadToJSON(p OutboxPayload) ([]byte, error) {
	return json.Marshal(p)
}
