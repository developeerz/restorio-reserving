# Структура проекта

cmd/                — точка входа (main.go)
internal/
  config/           — загрузка .env, флагов, конфигурации
    config.go
  db/               — инициализация базы (driver-адаптер)
    database.go
  kafka/            — адаптер Kafka (driver-адаптер)
    kafka.go
  dto/              — структуры для обмена данными (вход-выход HTTP, RPC)
    create_table_request.go
    reservation_request.go
    user_reservation_response.go
    time_slot_response.go
  mapper/           — преобразователи между сущностями/DTO/models
    payload_mapper.go
  domain/           — ваши «чистые» доменные сущности и бизнес-правила
    table.go
    reservation.go
  port/             — порты (интерфейсы) для репозиториев, событий и т.п.
    repository.go      // интерфейс TableRepository, ReservationRepository...
    notifier.go        // интерфейс KafkaNotifier
  usecase/          — «интеракторы»: сценарии использования
    table_usecase.go   // CreateTable, MoveTable, DeleteTable…
    reservation_usecase.go
  adapter/          — конкретные реализации портов
    postgres/
      table_repo.go           // PostgresTableRepository implements TableRepository
      reservation_repo.go
      outbox_repo.go
    kafka/
      kafka_notifier.go       // implements KafkaNotifier via internal/kafka
    http/
      handlers.go             // все HTTP-хэндлеры используют usecase
  pkg/              — модели для внешних систем (e.g. pkg/models для Kafka-пейлоадов)
docs/               — swagger, схемы, миграции и т.п.
