package service

import (
	"context"
	"excel-service/internal/models"
	"excel-service/internal/repository"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	container "github.com/vielendanke/go-db-lb"
	"github.com/xuri/excelize/v2"
)

type ExcelServiceImpl struct {
	repo repository.ExcelRepository
	lb   *container.LoadBalancer
}

func NewExcelService(repo repository.ExcelRepository, lb *container.LoadBalancer) ExcelService {
	return &ExcelServiceImpl{repo: repo, lb: lb}
}

func (e ExcelServiceImpl) SaveExcelFile(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
	src, err := file.Open()
	if err != nil {
		log.Errorf("failed ti open file: %v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	defer src.Close()

	excelFile, fileErr := excelize.OpenReader(src)
	if fileErr != nil {
		log.Errorf("failed to open reader: %v", fileErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fileErr)
	}

	// Get all the rows in the Sheet1.
	rows, rowsErr := excelFile.GetRows("Шаблон")
	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
	}

	var nomenclatures []*models.Nomenclature

	for i, row := range rows {
		if i == 0 || i == 1 {
			continue
		}
		nomenclature := &models.Nomenclature{}
		nomenclature.Id = uuid.New().String()
		nomenclature.PackageId = uuid.New().String()
		nomenclature.CodeSkmtr = row[0]
		nomenclature.CodeKsNsi = row[1]
		nomenclature.CodeAmto = row[2]
		nomenclature.OKPD2 = row[3]
		nomenclature.CodeTnved = row[4]
		nomenclature.Name = row[5]
		nomenclature.TmcCodeVendor = row[7]
		nomenclature.TmcMark = row[8]
		nomenclature.GostTu = row[9]
		nomenclature.DateOfManufacture = row[10]
		nomenclature.Manufacturer = row[11]
		nomenclature.BatchNumber = row[12]
		if row[13] == "облагается" {
			nomenclature.IsTax = true
		}
		taxPercentage, taxErr := strconv.ParseFloat(row[14], 32)
		if taxErr != nil {
			log.Errorf("failed to parse string to float: %v", taxErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат процента налога")
		}
		nomenclature.TaxPercentage = float32(taxPercentage)
		pricePerUnit, unitPriceErr := strconv.ParseFloat(row[15], 32)
		if unitPriceErr != nil {
			log.Errorf("failed to parse string to float: %v", unitPriceErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат цены за идиницу")
		}
		nomenclature.PricePerUnit = float32(pricePerUnit)
		nomenclature.Measurement = row[16]
		nomenclature.PriceValidThrough = row[17]
		wholesalePrice, wpErr := strconv.ParseFloat(row[18], 32)
		if wpErr != nil {
			log.Errorf("failed to parse string to float: %v", wpErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат оптовой цены за ед")
		}
		wholesaleItems := &models.WholesaleItems{}
		wholesaleItems.WholesalePricePerUnit = float32(wholesalePrice)
		wholesaleOrderFrom, woFromErr := strconv.Atoi(row[19])
		if woFromErr != nil {
			log.Errorf("failed to parse string to int: %v", woFromErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат оптовый заказ от")
		}
		wholesaleOrderTo, woToErr := strconv.Atoi(row[20])
		if woToErr != nil {
			log.Errorf("failed to parse string to int: %v", woToErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат оптовый заказ от")
		}
		wholesaleItems.WholesaleOrderFrom = wholesaleOrderFrom
		wholesaleItems.WholesaleOrderTo = wholesaleOrderTo
		nomenclature.WholesaleItems = wholesaleItems
		quantity, qErr := strconv.Atoi(row[21])
		if qErr != nil {
			log.Errorf("failed to parse string to int: %v", qErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат количество")
		}
		nomenclature.Quantity = quantity
		if row[22] == "в наличии" {
			nomenclature.ProductAvailability = true
		}
		nomenclature.HazardClass = row[26]
		nomenclature.PackagingType = row[27]
		nomenclature.PackingMaterial = row[28]
		nomenclature.StorageType = row[32]
		length, lenErr := strconv.ParseFloat(row[33], 32)
		if lenErr != nil {
			log.Errorf("float to parse length to float: %v", lenErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат длины")
		}
		width, widthErr := strconv.ParseFloat(row[34], 32)
		if widthErr != nil {
			log.Errorf("float to parse width to float: %v", widthErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат ширины")
		}
		height, heightErr := strconv.ParseFloat(row[35], 32)
		if heightErr != nil {
			log.Errorf("float to parse height to float: %v", heightErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат высоты")
		}
		nomenclature.Length = float32(length)
		nomenclature.Width = float32(width)
		nomenclature.Height = float32(height)
		amountInPackage, amountErr := strconv.ParseFloat(row[36], 32)
		if amountErr != nil {
			log.Errorf("failed to parse amount in package: %v", amountErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат количество в упаковке")
		}
		nomenclature.AmountInPackage = int8(amountInPackage)

		wNetto, wNettoErr := strconv.Atoi(row[37])
		if wNettoErr != nil {
			log.Errorf("failed to parse string to int: %v", wNettoErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(нетто)")
		}
		wBrutto, wBruttoErr := strconv.Atoi(row[38])
		if wBruttoErr != nil {
			log.Errorf("failed to parse string to int: %v", wBruttoErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(брутто)")
		}
		nomenclature.WeightNetto = int16(wNetto)
		nomenclature.WeightBrutto = int16(wBrutto)
		nomenclature.LoadingType = row[40]
		nomenclature.WarehouseAddress = row[41]
		nomenclature.Regions = row[42]
		nomenclature.DeliveryType = row[43]

		nomenclatures = append(nomenclatures, nomenclature)
	}

	tx, txErr := e.lb.CallPrimaryPreferred().PGxPool().Begin(ctx)
	if txErr != nil {
		log.Errorf("failed to begin tx: %v", txErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, txErr)
	}
	defer func(ctx context.Context) {
		cErr := tx.Commit(ctx)
		if cErr != nil {
			log.Errorf("failed to commit tx in service: %v", cErr)
			return
		}
	}(ctx)

	for _, v := range nomenclatures {
		repoErr := e.repo.SaveNomenclature(ctx, v, tx)
		if repoErr != nil {
			return nil, repoErr
		}
	}

	return &models.ResponseMsg{Message: "success"}, nil
}
