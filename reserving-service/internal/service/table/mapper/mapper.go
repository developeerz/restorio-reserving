package mapper

import (
	"github.com/developeerz/restorio-reserving/reserving-service/internal/dto"
	"github.com/developeerz/restorio-reserving/reserving-service/internal/repository/postgres/entity/table"
)

func mapToTable(table *table.Table) dto.Table {
	return dto.Table{
		TableID:     table.TableID,
		TableNumber: table.TableNumber,
		SeatsNumber: table.SeatsNumber,
		Type:        table.Type,
		Shape:       table.Shape,
		X:           table.X,
		Y:           table.Y,
	}
}

func MapToTables(tables []table.Table) []dto.Table {
	tablesDTO := make([]dto.Table, 0, len(tables))

	for i, table := range tables {
		tablesDTO[i] = mapToTable(&table)
	}

	return tablesDTO
}
