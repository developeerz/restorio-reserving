CREATE TABLE IF NOT EXISTS "Positions" (
	"table_id" INTEGER NOT NULL UNIQUE GENERATED BY DEFAULT AS IDENTITY,
	"x" INTEGER,
	"y" INTEGER,
	PRIMARY KEY("table_id"),
	FOREIGN KEY ("table_id") REFERENCES "Tables"("table_id") ON DELETE CASCADE
);