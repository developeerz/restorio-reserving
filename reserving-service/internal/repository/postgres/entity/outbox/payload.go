package outbox

type Payload struct {
	RestaurantName    string `db:"name"`
	RestaurantAddress string `db:"address"`
	TableNumber       string `db:"table_number"`
}
