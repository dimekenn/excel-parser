package handler

import (
	"excel-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type Handler struct {
	excelService service.ExcelService
}

func NewHandler(excelService service.ExcelService) *Handler {
	return &Handler{excelService: excelService}
}

func (h *Handler) SaveExcelFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.SaveExcelFile(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}
