package service

import (
	"context"
	"excel-service/internal/models"
	"mime/multipart"
)

type ExcelService interface {
	SaveExcelFile(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error)
}
