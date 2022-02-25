package repository

import (
	"context"
	"excel-service/internal/models"
	"github.com/jackc/pgx/v4"
)

type ExcelRepository interface {
	SaveNomenclature(ctx context.Context, nomenclature *models.Nomenclature, tx pgx.Tx) (string, error)
}
