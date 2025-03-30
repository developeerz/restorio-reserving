CREATE TABLE "restaurants" (
	"restaurant_id" INTEGER NOT NULL UNIQUE GENERATED BY DEFAULT AS IDENTITY,
	"name" VARCHAR(255) NOT NULL,
	"address" TEXT NOT NULL,
	"phone_number" VARCHAR(31),
	PRIMARY KEY("restaurant_id")
);