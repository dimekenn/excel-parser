package repository

import (
	"context"
	"excel-service/internal/models"

	"github.com/jackc/pgx/v4"
)

type ExcelRepository interface {
	SaveNomenclature(ctx context.Context, nomenclature *models.Nomenclature, tx pgx.Tx, userId, companyId string) error
	SaveArrayNomenclature(ctx context.Context, nomenclatures []*models.Nomenclature, tx pgx.Tx) error
	SaveMTRFile(ctx context.Context, nomenclature *models.Mtr, tx pgx.Tx) error
	NewParentCategory(ctx context.Context, cat string, tx pgx.Tx) error
	NewChildCategory(ctx context.Context, cat *models.Category, tx pgx.Tx) error
	CheckCategory(ctx context.Context, catName string, tx pgx.Tx) (bool, error)
	CheckCompany(ctx context.Context, inn string) (bool, error)
	CreateCompany(ctx context.Context, company *models.Company, tx pgx.Tx) error
	CreateUserByCompany(ctx context.Context, inn, email, companyId, companyName string) error
	SelectUser(ctx context.Context, inn string) (string, error)
	SelectCompanyInnById(ctx context.Context, companyId string) (string, error)
	SelectPriceListsByUploadId(ctx context.Context, uploadId string) ([]string, error)
	SetUploadStatus(ctx context.Context, uploadId string, status string) error
	SaveBanks(ctx context.Context, bik, name, cor_account, address string, tx pgx.Tx) error
	NewErrorNomenclatureId(ctx context.Context, row_id int, fileName string) error
	NewUploadCatalogue(ctx context.Context, fileNameDisc, fileNameDl, uploadedBy, companyId string, fileSize int64) error
	GetFromUploadCatalogue(ctx context.Context, id string) ([]*models.UploadsEntity, error)
}
