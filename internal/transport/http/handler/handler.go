package handler

import (
	"excel-service/internal/models"
	"excel-service/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	excelService service.ExcelService
}

func NewHandler(excelService service.ExcelService) *Handler {
	return &Handler{excelService: excelService}
}

//SaveExcelFile godoc
//@Summary Parsing excel file from supplier
//@Description accept multipart/form-data returns json struct
//@Accept mpfd
//@Produce json
//@Param excel file formData file true "file"
//@Success 200 {object} models.ResponseMsg
//@Failure 400 {object} models.ResponseMsg
//@Failure 500 {object} models.ResponseMsg
//@Router /api/v1/upload/excel [post]
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

func (h *Handler) SaveMtr(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.SaveMTRExcelFile(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) NewCategory(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.SaveCategory(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) NewCompany(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.CreateCompany(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SaveOrganizerNomenclature(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.SaveOrganizerNomenclature(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SaveBanks(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	res, resErr := h.excelService.SaveBanks(c.Request().Context(), file)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetExcelFromAwsByFileId(c echo.Context) error {
	var req models.GetExcelFromAwsByFileIdReq
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}

	res, resErr := h.excelService.GetExcelFromAwsByFileId(c.Request().Context(), &req)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ProcessDirectusUploads(c echo.Context) error {
	var req models.DirectusModel
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}

	res, resErr := h.excelService.SaveNomenclatureFromDirectus(c.Request().Context(), &req)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UploadExcelFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read file by key file")
	}

	companyName := c.FormValue("company_name")
	if companyName == "" {
		log.Errorf("failed to read company name")
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read company name by key company_name")
	}

	res, resErr := h.excelService.UploadExcelFile(c.Request().Context(), file, companyName)
	if resErr != nil {
		return resErr
	}

	log.Infof("success response: %v", res)
	return c.JSON(http.StatusOK, res)
}
