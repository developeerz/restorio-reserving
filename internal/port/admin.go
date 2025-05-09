package port

import "github.com/developeerz/restorio-reserving/internal/dto"

// AdminRepository описывает, что нужно для админских операций
type AdminRepository interface {
	// CreateTable сохраняет новый столик и возвращает его ID
	CreateTable(req dto.CreateTableRequest) (int, error)
}
