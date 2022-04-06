package service

import (
	"context"
	"excel-service/internal/models"
	"mime/multipart"
)

type ExcelService interface {
	SaveExcelFile(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	SaveMTRExcelFile(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	SaveCategory(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	CreateCompany(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	SaveOrganizerNomenclature(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	SaveBanks(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
	GetExcelFromAwsByFileId(ctx context.Context, req *models.GetExcelFromAwsByFileIdReq) (*models.ResponseMsg, error)
	UploadExcelFile(ctx context.Context, file *multipart.FileHeader, companuyName string) (*models.ResponseMsg, error)
	SaveNomenclatureFromDirectus(ctx context.Context, req *models.DirectusModel) (*models.ResponseMsg, error)
}
