CREATE TABLE IF NOT EXISTS "tables" (
	"table_id" INTEGER NOT NULL UNIQUE GENERATED BY DEFAULT AS IDENTITY,
	"restaurant_id" INTEGER NOT NULL,
	"table_number" VARCHAR(255),
	"seats_number" INTEGER NOT NULL,
	"type" VARCHAR(63),
	PRIMARY KEY("table_id"),
	FOREIGN KEY ("restaurant_id") REFERENCES "restaurants"("restaurant_id") ON DELETE CASCADE
);