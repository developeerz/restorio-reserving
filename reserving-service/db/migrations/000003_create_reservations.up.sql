CREATE TABLE IF NOT EXISTS "reservations" (
	"reservation_id" INTEGER NOT NULL UNIQUE GENERATED BY DEFAULT AS IDENTITY,
	"user_id" INTEGER NOT NULL,
	"table_id" INTEGER NOT NULL,
	"reservation_time_from" TIMESTAMP NOT NULL,
	"reservation_time_to" TIMESTAMP NOT NULL,
	"status" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("reservation_id"),
	FOREIGN KEY ("table_id") REFERENCES "tables"("table_id") ON DELETE CASCADE
);