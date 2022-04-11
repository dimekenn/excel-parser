package service

import (
	"context"
	"excel-service/internal/configs"
	"excel-service/internal/models"
	"excel-service/internal/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	container "github.com/vielendanke/go-db-lb"
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ExcelServiceImpl struct {
	repo repository.ExcelRepository
	lb   *container.LoadBalancer
	cfg  *configs.Configs
}

func NewExcelService(repo repository.ExcelRepository, lb *container.LoadBalancer, cfg *configs.Configs) ExcelService {
	return &ExcelServiceImpl{repo: repo, lb: lb, cfg: cfg}
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

	saveErr := NewMTRFile(rows, e.repo, ctx, "", "")
	if saveErr != nil {
		return nil, saveErr
	}

	return &models.ResponseMsg{Message: "success"}, nil
}

func newSupplierNomenclature(rows [][]string, priceLists []string, repo repository.ExcelRepository, ctx context.Context, companyId, userId string) error {
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
		if len(row) > 10 {
			nomenclature.DateOfManufacture = row[10]
		}
		if len(row) > 11 {
			nomenclature.Manufacturer = row[11]
		}
		if len(row) > 12 {
			nomenclature.BatchNumber = row[12]
		}
		//nomenclature.CompanyInn = companyInn
		if len(row) > 13 && row[13] == "облагается" {
			nomenclature.IsTax = true
		}
		if len(row) > 14 && nomenclature.IsTax {
			taxPercentage, taxErr := strconv.ParseFloat(row[14], 8)
			if taxErr != nil {
				log.Errorf("failed to parse string to float: %v", taxErr)
			} else {
				nomenclature.TaxPercentage = float32(taxPercentage)
			}
		}

		if len(priceLists) > 0 {
			nomenclature.PriceLists = priceLists
		}

		if len(row) > 16 {
			pricePerUnit, unitPriceErr := strconv.ParseFloat(row[15], 8)
			if unitPriceErr != nil {
				log.Errorf("failed to parse string to float: %v", unitPriceErr)
			} else {
				nomenclature.PricePerUnit = float32(pricePerUnit)
			}
		}

		if len(row) > 16 {
			nomenclature.Measurement = row[16]
		}
		if len(row) > 17 {
			nomenclature.PriceValidThrough = row[17]
		}

		wholesaleItems := &models.WholesaleItems{}

		if len(row) > 18 {
			wholesaleItems.WholesalePricePerUnit = row[18]
		}
		if len(row) > 19 && row[19] != "" {
			orderDateArr := strings.Split(row[19], "и")
			if len(orderDateArr) == 2 {
				// wholesaleOrderFrom, woFromErr := strconv.Atoi(orderDateArr[0])
				// if woFromErr != nil {
				// 	log.Errorf("failed to parse string to int: %v", woFromErr)
				// }
				wholesaleItems.WholesaleOrderFrom = orderDateArr[0]

				wholesaleItems.WholesaleOrderTo = orderDateArr[1]
			}

		}
		nomenclature.WholesaleItems = wholesaleItems

		if len(row) > 20 && row[20] != "" {
			quantity, qErr := strconv.Atoi(row[20])
			if qErr != nil {
				log.Errorf("failed to parse string to int: %v", qErr)
				return echo.NewHTTPError(http.StatusBadRequest, "не правильный формат количество")
			}

			nomenclature.Quantity = quantity
		}

		if len(row) > 21 && (row[21] == "в наличии" || row[21] == "да") {
			nomenclature.ProductAvailability = true
		}

		if len(row) > 25 {
			nomenclature.HazardClass = row[25]
		}
		if len(row) > 26 {
			nomenclature.PackagingType = row[26]
		}
		if len(row) > 27 {
			nomenclature.PackingMaterial = row[27]
		}
		if len(row) > 32 {
			nomenclature.StorageType = row[32]
		}
		if len(row) > 32 && row[32] != "" {
			length, lenErr := strconv.ParseFloat(row[32], 32)
			if lenErr != nil {
				log.Errorf("float to parse length to float: %v", lenErr)
			}
			nomenclature.Length = float32(length)

		}
		if len(row) > 33 && row[33] != "" {
			width, widthErr := strconv.ParseFloat(row[33], 32)
			if widthErr != nil {
				log.Errorf("float to parse width to float: %v", widthErr)
			}
			nomenclature.Width = float32(width)

		}

		if len(row) > 34 && row[34] != "" {
			height, heightErr := strconv.ParseFloat(row[34], 32)
			if heightErr != nil {
				log.Errorf("float to parse height to float: %v", heightErr)
			}
			nomenclature.Height = float32(height)
		}

		if len(row) > 35 && row[35] != "" {
			amountInPackage, amountErr := strconv.Atoi(row[35])
			if amountErr != nil {
				log.Errorf("failed to parse amount in package: %v", amountErr)
			}
			nomenclature.AmountInPackage = int8(amountInPackage)
		}

		if len(row) > 36 && row[36] != "" {
			wNetto, wNettoErr := strconv.Atoi(row[36])
			if wNettoErr != nil {
				log.Errorf("failed to parse string to int: %v", wNettoErr)
			}
			nomenclature.WeightNetto = float32(wNetto)

		}

		if len(row) > 37 && row[37] != "" {
			wBrutto, wBruttoErr := strconv.Atoi(row[37])
			if wBruttoErr != nil {
				log.Errorf("failed to parse string to int: %v", wBruttoErr)
			}
			nomenclature.WeightBrutto = float32(wBrutto)
		}

		if len(row) > 38 && row[38] != "" {
			volume, volumeErr := strconv.Atoi(row[38])
			if volumeErr != nil {
				log.Errorf("failed to parse string to int: %v", volumeErr)
			}
			nomenclature.Volume = float32(volume)
		}

		if len(row) > 39 {
			nomenclature.LoadingType = row[39]
		}
		if len(row) > 40 {
			nomenclature.WarehouseAddress = row[40]
		}
		if len(row) > 41 {
			nomenclature.Regions = row[41]
		}
		if len(row) > 42 {
			nomenclature.DeliveryType = row[42]
		}

		saveErr := repo.SaveNomenclature(ctx, nomenclature, nil, userId, companyId)
		if saveErr != nil {
			repo.NewErrorNomenclatureId(ctx, i, "supplier_nomenclature")
		}
	}
	return nil

}

func NewMTRFile(rows [][]string, repo repository.ExcelRepository, ctx context.Context, userId, companyId string) error {
	for i, v := range rows {
		fmt.Println("started")
		if i == 0 {
			continue
		}
		nomenclature := &models.Nomenclature{}
		nomenclature.Id = uuid.New().String()
		nomenclature.PackageId = uuid.New().String()
		name := v[5]
		if v[6] != "" {
			name = strings.Replace(name, "("+v[6]+")", "", 1)
		}
		name, nomenclature.Length, nomenclature.Height, nomenclature.Width = takeVolume(name)

		if len(v[28]) < 2 {
			name, nomenclature.DrawingName = takeDraw(name)
		} else {
			nomenclature.DrawingName = v[28]
		}
		fmt.Println("drawing name: ", nomenclature.DrawingName)
		nomenclature.Name = name
		nomenclature.CodeSkmtr = v[6]
		nomenclature.Measurement = v[7]
		nomenclature.OKPD2 = v[15]
		nomenclature.CodeTnved = v[20]
		nomenclature.CodeAmto = v[25]
		if v[29] != "" {
			wNetto, wNettoErr := strconv.ParseFloat(v[29], 8)
			if wNettoErr != nil {
				log.Errorf("failed to parse string to int: %v", wNettoErr)
				return echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(нетто): "+v[29])
			}
			nomenclature.WeightNetto = float32(wNetto)
		}

		if v[30] != "" {
			wBrutto, wBruttoErr := strconv.ParseFloat(v[30], 8)
			if wBruttoErr != nil {
				log.Errorf("failed to parse string to int: %v", wBruttoErr)
				return echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(брутто)"+v[30])
			}
			nomenclature.WeightBrutto = float32(wBrutto)
		} else {
			name, nomenclature.WeightBrutto = takeWeight(name, v[30])
		}

		name, nomenclature.WeightNetto = takeWeight(name, v[29])

		nomenclature.Name = name

		nomenclature.TmcMark = v[41]
		nomenclature.Manufacturer = v[19]
		nomenclature.FullName = v[17]
		nomenclature.Representation = v[52]
		nomenclature.Link = v[0]

		if v[10] == "" {
			nomenclature.CategoryName = v[10]
		}
		nomenclature.CategoryName = "He классифицированные"
		//nomenclature.DeliveryAddress =

		nomenclatureMTR := &models.Mtr{}
		wholesaleItems := &models.WholesaleItems{}

		//nomenclatureMTR.Link = v[0]
		nomenclatureMTR.DataVersion = v[2]
		nomenclatureMTR.DeleteMark = v[3]
		nomenclatureMTR.Code = v[4]
		//nomenclatureMTR.Name = v[5]
		//nomenclatureMTR.VendorCode = v[6]
		//nomenclatureMTR.Measurement = v[7]
		nomenclatureMTR.Identifier = v[8]
		nomenclatureMTR.CatalogueNumber = v[9]
		//nomenclatureMTR.Class = v[10]
		nomenclatureMTR.Comments = v[11]
		nomenclatureMTR.PropertySet = v[12]
		nomenclatureMTR.TechDoc = v[13]
		nomenclatureMTR.Okved2 = v[14]
		//nomenclatureMTR.Okpd2 = v[15]
		nomenclatureMTR.Description = v[16]
		nomenclatureMTR.FullName = v[17]
		nomenclatureMTR.SignOfUser = v[18]
		//nomenclatureMTR.Manufacturer = v[19]
		//nomenclatureMTR.Tnved = v[20]
		nomenclatureMTR.DeleteRecord = v[21]
		nomenclatureMTR.DeleteItemType = v[22]
		nomenclatureMTR.DeleteRefPosition = v[23]
		nomenclatureMTR.DeleteLayout = v[24]
		//nomenclatureMTR.SlAmto = v[25]
		nomenclatureMTR.SlManufacturerVendorCode = v[26]
		nomenclatureMTR.SlManufacturerBarcode = v[27]
		//nomenclatureMTR.SlDraw = v[28]
		//nomenclatureMTR.SlWeightNetto = v[29]
		//nomenclatureMTR.SlWeightBrutto = v[30]
		nomenclatureMTR.SlPriority = v[31]
		nomenclatureMTR.SlSupplierMeasurement = v[32]
		nomenclatureMTR.SlConversionFactor = v[33]
		nomenclatureMTR.SlSupplierWeightNetto = v[34]
		nomenclatureMTR.SlSupplierWeightBrutto = v[35]
		nomenclatureMTR.SlExpiryDate = v[36]
		nomenclatureMTR.SlManufacturerCountry = v[37]
		nomenclatureMTR.SlCheckInterval = v[38]
		nomenclatureMTR.SlDrawingFile = v[39]
		nomenclatureMTR.SlImgFile = v[40]
		//nomenclatureMTR.SlMarkTmc = v[41]
		nomenclatureMTR.SlStateStandard = v[42]
		nomenclatureMTR.SlPackage = v[43]
		nomenclatureMTR.SlHazardClass = v[44]
		nomenclatureMTR.SlNomenclatureSign = v[45]
		nomenclatureMTR.SlSize = v[46]
		nomenclatureMTR.MdmKey = v[47]
		nomenclatureMTR.NsiRequest = v[48]
		nomenclatureMTR.NsiManualChange = v[49]
		nomenclatureMTR.Predefined = v[50]
		nomenclatureMTR.PredefinedDataName = v[51]
		//nomenclatureMTR.Representation = v[52]
		//nomenclatureMTR.Measurement1 = v[53]
		//nomenclatureMTR.Coefficient = v[54]
		//nomenclatureMTR.Purpose = v[55]
		//nomenclatureMTR.Analog = v[56]
		//nomenclatureMTR.KindOfClassifier = v[57]
		//nomenclatureMTR.Class1 = v[58]
		//nomenclatureMTR.Property = v[59]
		//nomenclatureMTR.Value = v[60]
		//nomenclatureMTR.TextString = v[61]
		//nomenclatureMTR.SparePart = v[62]
		//nomenclatureMTR.Shipper = v[63]
		//nomenclatureMTR.ShippingAddress = v[64]
		//nomenclatureMTR.MinShippingBatch = v[65]
		//nomenclatureMTR.CharacteristicName = v[66]
		//nomenclatureMTR.Characteristic = v[67]
		//nomenclatureMTR.Value1 = v[68]

		nomenclature.Payload = nomenclatureMTR
		nomenclature.WholesaleItems = wholesaleItems

		err := repo.SaveNomenclature(ctx, nomenclature, nil, userId, companyId)
		if err != nil {
			repo.NewErrorNomenclatureId(ctx, i, "mtr")
		}
	}
	return nil
}

func (e ExcelServiceImpl) SaveMTRExcelFile(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
	src, err := file.Open()
	if err != nil {
		log.Errorf("failed ti open file: %v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	defer src.Close()
	fmt.Println("start")
	excelFile, fileErr := excelize.OpenReader(src)
	if fileErr != nil {
		log.Errorf("failed to open reader: %v", fileErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fileErr)
	}

	// Get all the rows in the Sheet1.
	rows, rowsErr := excelFile.GetRows("TDSheet")
	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
	}

	newMtrErr := NewMTRFile(rows, e.repo, ctx, "", "")
	if newMtrErr != nil {
		return nil, newMtrErr
	}

	return &models.ResponseMsg{Message: "success"}, nil
}


func takeDraw(name string) (string, string) {
	fmt.Println("takeDraw")
	regDraw, dErr := regexp.Compile(`[а-яА-Я0-9]{1,4}([.][а-яА-Я0-9]{1,4}){2,5}`)
	if dErr != nil {
		log.Errorf("failed to regexp compile in takeDraw: %v", dErr)
	}

	names := regDraw.FindAllString(name, 5)
	if len(names) == 1 {
		draw := names[0]
		name = strings.Replace(name, draw, "", 1)
		fmt.Println("draw", draw)
		return name, draw
	} else if len(names) > 1 {
		draw := names[len(names)-1]
		for _, v := range names {
			if v == draw {
				name = strings.Replace(name, draw, "", 1)
			}
		}
		name = strings.Replace(name, draw, "", 1)
		fmt.Println("draw", draw)
		return name, draw
	}

	return name, ""
}

func takeWeight(name, weightCell string) (string, float32) {
	regWei, wErr := regexp.Compile(`[0-9]{1,4}[,]*[0-9]{1,5}кг`)
	if wErr != nil {
		log.Errorf("failed to regexp compile in takeWeight: %v", wErr)
	}
	weight := regWei.FindString(name)
	fmt.Println("weight", weight)
	name = strings.Replace(name, weight, "", 1)
	if len(weight) > 3 {
		name = strings.Replace(name, weight, "", 1)
		weight = strings.Replace(weight, "кг", "", 1)
		weight = strings.Replace(weight, ",", ".", 1)
		weiFloat, weErr := strconv.ParseFloat(weight, 8)
		if weErr != nil {
			log.Errorf("failed to parse float in takeWeight: %v", weErr)
		}
		return name, float32(weiFloat)
	}
	return name, 0
}

func (e ExcelServiceImpl) SaveCategory(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
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
	fmt.Println(excelFile.GetSheetList())
	// Get all the rows in the Sheet1.
	fmt.Println("1")
	rows, rowsErr := excelFile.GetRows("TDSheet")
	fmt.Println("1")
	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
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

	catMap := make(map[string]bool)

	fmt.Println("bасталды")
	for i, v := range rows {
		if i == 0 {
			continue
		}
		if !catMap[v[10]] {

			if len(v[10]) < 1 {
				continue
			}

			// catArr := strings.Split(v[10], ", ")
			// if len(catArr) > 0 {
			// 	for _, a := range catArr {
			// 		isExist, isErr := e.repo.CheckCategory(ctx, a, tx)
			// 		if isErr != nil {
			// 			return nil, isErr
			// 		}
			// 		if isExist {
			// 			continue
			// 		}
			// 		err := e.repo.NewParentCategory(ctx, a, tx)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		catMap[a] = true
			// 	}
			// 	continue
			// }

			isExist, isErr := e.repo.CheckCategory(ctx, v[10], tx)
			if isErr != nil {
				return nil, isErr
			}
			if isExist {
				continue
			}

			err := e.repo.NewParentCategory(ctx, v[10], tx)
			if err != nil {
				return nil, err
			}

			catMap[v[10]] = true
		} else {
			continue
		}
	}

	//tx, txErr := e.lb.CallPrimaryPreferred().PGxPool().Begin(ctx)
	//if txErr != nil {
	//	log.Errorf("failed to begin tx: %v", txErr)
	//	return nil, echo.NewHTTPError(http.StatusInternalServerError, txErr)
	//}
	//defer func(ctx context.Context) {
	//	cErr := tx.Commit(ctx)
	//	if cErr != nil {
	//		log.Errorf("failed to commit tx in service: %v", cErr)
	//		return
	//	}
	//}(ctx)

	//var childCategories *[]models.Category

	//for _, v := range rows{
	//	category := &models.Category{}
	//	if len(v[0]) > 3{
	//		category.Name = v[0][3:]
	//	}
	//	category.Name = v[2]
	//	category.Code = v[3]
	//	category.ParentName = v[0]
	//	err := e.repo.NewChildCategory(ctx, category, tx)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return &models.ResponseMsg{Message: "success"}, nil
}

func (e ExcelServiceImpl) CreateCompany(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
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
	fmt.Println(excelFile.GetSheetList())
	// Get all the rows in the Sheet1.
	fmt.Println("1")
	rows, rowsErr := excelFile.GetRows("Лист1")
	fmt.Println("1")
	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
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

	comMap := make(map[string]bool)

	for i, row := range rows {
		if i < 10 {
			continue
		}
		if comMap[row[10]] {
			continue
		}
		comOk, checkErr := e.repo.CheckCompany(ctx, row[10])
		if checkErr != nil {
			return nil, checkErr
		}
		if comOk {
			continue
		}

		company := &models.Company{}
		company.Name = row[11]
		company.Inn = row[10]
		company.UserId = "d3162f03-6c63-42be-b31c-22542245074f"
		createErr := e.repo.CreateCompany(ctx, company, tx)
		if createErr != nil {
			fmt.Println(createErr)
		}
		comMap[row[10]] = true
		fmt.Println("inserted: ", company.Name)
	}
	return &models.ResponseMsg{Message: "success"}, nil
}

func (e ExcelServiceImpl) SaveOrganizerNomenclature(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
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
	fmt.Println(excelFile.GetSheetList())
	// Get all the rows in the Sheet1.
	fmt.Println("1")
	rows, rowsErr := excelFile.GetRows("Лист1")
	fmt.Println("1")
	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
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

	if orgRepErr := newOrgranizerNomenclature(rows, e.repo, ctx, "", ""); orgRepErr != nil {
		return nil, err
	}
	// saveErr := e.repo.SaveArrayNomenclature(ctx, nomenclatures, tx)
	// if saveErr != nil {
	// 	fmt.Println("save array nom errors: ", saveErr)
	// 	return nil, saveErr
	// }
	return &models.ResponseMsg{Message: "success"}, nil
}

func newOrgranizerNomenclature(rows [][]string, repo repository.ExcelRepository, ctx context.Context, userId, companyId string) error {
	for i, row := range rows {
		if i < 1 {
			continue
		}
		fmt.Println("row #", i)
		nomenclature := &models.Nomenclature{}
		nomenclature.Id = uuid.New().String()
		nomenclature.Name = row[65]
		nomenclature.TmcCodeVendor = row[14]
		nomenclature.Manufacturer = row[42]
		nomenclature.CodeTnved = row[53]
		nomenclature.CodeAmto = row[62]
		nomenclature.GostTu = row[69]
		nomenclature.DrawingName = row[70]
		nomenclature.CodeKsNsi = row[85]
		nomenclature.OKPD2 = row[90]
		nomenclature.CodeSkmtr = row[93]
		nomenclature.FullName = row[100]
		nomenclature.Representation = row[116]
		nomenclature.Measurement = row[24]
		nomenclature.Link = row[13]
		if len(row[27]) > 1 {
			taxPercentageString := strings.Replace(row[27], "%", "", 1)
			taxPercentage, taxErr := strconv.ParseFloat(taxPercentageString, 8)
			if taxErr != nil {
				log.Errorf("failed to parse tax percentage: %v", taxErr)
			} else {
				nomenclature.TaxPercentage = float32(taxPercentage)
				nomenclature.IsTax = true
			}
		}

		if len(row[11]) > 3 {
			nomenclature.UserId = row[11]
		} else {
			nomenclature.UserId = "Organizer"
		}

		// if len(row[10]) > 9 {
		// 	userId, idEr := e.repo.SelectUser(ctx, row[10])
		// 	if idEr != nil {
		// 		log.Errorf("%v", idEr)
		// 	}

		// 	nomenclature.UserId = row[10]

		// 	nomenclature.CompanyInn = row[10]
		// }

		orgNomenclature := &models.OrganizerNomenclature{}
		orgNomenclature.NomenclatureType = row[15]
		orgNomenclature.IsWeight = netFunc(row[16])
		orgNomenclature.WeightCoefficient = row[17]
		orgNomenclature.WIPBalance = row[18]
		orgNomenclature.PartitionAccountingBySeries = row[19]
		orgNomenclature.AccountingBySeries = row[20]
		orgNomenclature.KeepAccountingBySeriesWCD = row[21]
		orgNomenclature.KeepAccountingAccordingToCharacteristics = row[22]
		orgNomenclature.MainMeasurement = row[24]
		orgNomenclature.ReportMeasurement = row[25]
		orgNomenclature.ResidueMeasurement = row[26]
		orgNomenclature.Kit = row[28]
		orgNomenclature.PurposeOfUse = row[29]
		orgNomenclature.Comments = row[30]
		orgNomenclature.Service = row[31]
		orgNomenclature.NomenclatureGroup = row[33]
		orgNomenclature.FileImg = row[34]
		orgNomenclature.MainSupplier = row[35]
		orgNomenclature.SalesManager = row[36]
		orgNomenclature.ManufacturerCountry = row[37]
		orgNomenclature.GTDNumber = row[38]
		orgNomenclature.ArticleCost = row[39]
		orgNomenclature.RequiresExternalCertification = netFunc(row[40])
		orgNomenclature.RequiresInternalCertification = netFunc(row[41])
		orgNomenclature.Set = netFunc(row[44])
		orgNomenclature.OKP = row[49]
		orgNomenclature.IsAlcohol = netFunc(row[55])
		orgNomenclature.IsImportAlcohol = netFunc(row[56])
		orgNomenclature.VolumeDAL = row[58]
		orgNomenclature.QuarantineZone = netFunc(row[60])
		orgNomenclature.CodeSUMI = row[61]
		orgNomenclature.AMTOStatus = row[63]
		orgNomenclature.ENSKStatus = row[64]
		orgNomenclature.ENSKName = row[65]
		orgNomenclature.ENSKTM = row[66]
		orgNomenclature.ENSKBrandDesign = row[67]
		orgNomenclature.ENSKTechSpec = row[68]
		orgNomenclature.ENSKMaterialMark = row[71]
		orgNomenclature.ENSKGostMaterial = row[72]
		orgNomenclature.CatalogueNumber = row[73]
		orgNomenclature.ENSKOKPClassificator = row[74]
		orgNomenclature.AMTONormName = row[75]
		orgNomenclature.AMTOCodeForEOrder = row[76]
		orgNomenclature.ENSKExpertComments = row[77]
		orgNomenclature.TMXClassificatorGP = row[78]
		orgNomenclature.TMXClassificatorOKP = row[79]
		orgNomenclature.TMXClassificatorRTK = row[80]
		orgNomenclature.TMXCodePDM = row[94]
		orgNomenclature.TMXItemType = row[95]
		orgNomenclature.IsTobacco = netFunc(row[101])
		orgNomenclature.IsShoes = netFunc(row[103])
		orgNomenclature.TMXCodeMDM = row[108]

		nomenclature.OrganizerNomenclature = orgNomenclature
		//nomenclatures = append(nomenclatures, nomenclature)
		err := repo.SaveNomenclature(ctx, nomenclature, nil, userId, companyId)
		if err != nil {
			log.Error(err)
			repo.NewErrorNomenclatureId(ctx, i, "organizer_nomenclature")
		}
	}
	return nil
}

func netFunc(is string) bool {
	if is == "нет" {
		return false
	}
	return true
}

func (e ExcelServiceImpl) SaveBanks(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
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

	fmt.Println(excelFile.GetSheetList())
	// Get all the rows in the Sheet1.

	rows, rowsErr := excelFile.GetRows("Лист1")

	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
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

	//var nomenclatures []*models.Nomenclature

	for i, row := range rows {
		if i < 4 {
			continue
		}
		if row[9] == "Банк" {
			if len(row[7]) < 2 {
				continue
			}
			bik := row[0]
			addres := row[3] + ", " + row[4]
			name := row[5]
			cor_account := row[7]
			err := e.repo.SaveBanks(ctx, bik, name, cor_account, addres, tx)
			if err != nil {
				log.Errorf("save bank error: %v", err)
				return nil, err
			}
			//return &models.ResponseMsg{Message: "success"}, nil
		}

	}
	return &models.ResponseMsg{Message: "success"}, nil
}

func (e ExcelServiceImpl) GetExcelFromAwsByFileId(ctx context.Context, req *models.GetExcelFromAwsByFileIdReq) (*models.ResponseMsg, error) {
	endpoint := e.cfg.Aws.Host
	accessKeyID := e.cfg.Aws.SecretKey
	secretAccessKey := e.cfg.Aws.AccessKey
	bucket := e.cfg.Aws.Bucket

	useSSL := false

	log.Error(endpoint)
	log.Error("db", e.cfg.DB.DBName)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(secretAccessKey, accessKeyID, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error("failed to connect to minio: ", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//filePath := fmt.Sprintf("%s", req.FileId)

	obj, err := minioClient.GetObject(ctx, bucket, req.FileId, minio.GetObjectOptions{})
	if err != nil {
		log.Error("get object err:", err)
		return nil, err
	}
	if obj == nil {
		log.Warn("object is nil")
		return nil, nil
	}

	return &models.ResponseMsg{Message: "success"}, nil
}

func (e ExcelServiceImpl) UploadExcelFile(ctx context.Context, file *multipart.FileHeader, companyName string) (*models.ResponseMsg, error) {
	src, err := file.Open()
	if err != nil {
		log.Errorf("failed ti open file: %v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	endpoint := e.cfg.Aws.Host
	accessKeyID := e.cfg.Aws.SecretKey
	secretAccessKey := e.cfg.Aws.AccessKey
	bucket := e.cfg.Aws.Bucket
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(secretAccessKey, accessKeyID, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error("failed to connect to minio: ", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	contentType := "application/vnd.ms-excel"
	fileNameDisc := uuid.New().String() + file.Filename[len(file.Filename)-5:]

	uploadInfo, uploadErr := minioClient.PutObject(ctx, bucket, fileNameDisc, src, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if uploadErr != nil {
		log.Error("failed to upload file to s3:", uploadErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, uploadErr)
	}

	log.Infof("file %s success uploaded, size: %d", file.Filename, uploadInfo.Size)

	err = e.repo.NewUploadCatalogue(ctx, fileNameDisc, file.Filename, "", companyName, file.Size)
	if err != nil {
		return nil, err
	}

	return &models.ResponseMsg{Message: "success"}, nil
}

func (e ExcelServiceImpl) SaveNomenclatureFromDirectus(ctx context.Context, req *models.DirectusModel) (*models.ResponseMsg, error) {
	if req.Collection != "uploads" {
		log.Warnf("collection is %s not uploads", req.Collection)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "collection is not uploads")
	}
	log.Infof("Start processing upload from directus: %v", req)
	err := e.repo.SetUploadStatus(ctx, req.Key, "processing")
	if err != nil {
		log.Warnf("failed to set upload status", err)
		return nil, err
	}
	//time.Sleep(15 * time.Second) // todo удалить после демо 8.10
	uploads, uploadErr := e.repo.GetFromUploadCatalogue(ctx, req.Key)
	if uploadErr != nil {
		log.Warnf("failed to get upload catalog", uploadErr)
		return nil, uploadErr
	}

	endpoint := e.cfg.Aws.Host
	accessKeyID := e.cfg.Aws.SecretKey
	secretAccessKey := e.cfg.Aws.AccessKey
	bucket := e.cfg.Aws.Bucket
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(secretAccessKey, accessKeyID, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error("failed to connect to minio: ", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, upload := range uploads{
		_, err := processFiles(minioClient, ctx, bucket, upload, req, e.repo)
		if err != nil{
			return nil, err
		}
	}
	// err = e.repo.SetUploadStatus(ctx, req.Key, "processed")
	// if err != nil {
	// 	return nil, err
	// }

	log.Errorf("failed to find correct template")
	return nil, echo.NewHTTPError(http.StatusBadRequest, "неправильный шаблон документа, обратитесь к нам")
}

func processFiles(minioClient *minio.Client, ctx context.Context, bucket string, upload *models.UploadsEntity, req *models.DirectusModel, repo repository.ExcelRepository) (*models.ResponseMsg, error){
	minioObj, getObjErr := minioClient.GetObject(ctx, bucket, upload.FileId, minio.GetObjectOptions{})
	if getObjErr != nil {
		log.Error("failed to get object:", getObjErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, upload.FileId+"does not exists")
	}

	defer minioObj.Close()

	excelFile, fileErr := excelize.OpenReader(minioObj)
	if fileErr != nil {
		log.Errorf("failed to open reader: %v", fileErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fileErr)
	}

	rows, rowsErr := excelFile.GetRows(excelFile.GetSheetList()[0])

	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
	}

	if rows[0][10] == "ИНН" && rows[0][11] == "Поставщик" {
		orgNomErr := newOrgranizerNomenclature(rows, repo, ctx, upload.UserId, upload.CompanyId)
		if orgNomErr != nil {
			log.Errorf("failed parse: %v", orgNomErr)
			return nil, orgNomErr
		}
		return &models.ResponseMsg{Message: "success"}, nil

	} else if rows[0][6] == "Наименование" && rows[0][7] == "Артикул" && rows[0][8] == "Идентификатор" {

		mtrErr := NewMTRFile(rows, repo, ctx, upload.UserId, upload.CompanyId)
		if mtrErr != nil {
			log.Errorf("failed parse: %v", mtrErr)
			return nil, mtrErr
		}
		err := repo.SetUploadStatus(ctx, req.Key, "processed")
		if err != nil {
			log.Errorf("failed to set upload status: %v", fileErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
		}
		return &models.ResponseMsg{Message: "success"}, nil

	} else if rows[0][0] == "Код СКМТР" && rows[0][1] == "КОД КС НСИ" && rows[0][2] == "Код АМТО" {

		// inn, companyErr := e.repo.SelectCompanyInnById(ctx, req.Accounting.Company)
		// if companyErr != nil {
		// 	log.Errorf("failed to get company inn: %v", companyErr)
		// 	return nil, echo.NewHTTPError(http.StatusBadRequest, "Внутренняя ошибка")
		// }

		priceLists, priceListErr := repo.SelectPriceListsByUploadId(ctx, req.Key)
		if priceListErr != nil {
			log.Errorf("failed to get price lists: %v", priceListErr)
			return  nil, echo.NewHTTPError(http.StatusBadRequest, "Внутренняя ошибка")
		}

		suppErr := newSupplierNomenclature(rows, priceLists, repo, ctx, upload.CompanyId, upload.UserId)
		if suppErr != nil {
			log.Errorf("failed parse: %v", suppErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
		}
		err := repo.SetUploadStatus(ctx, req.Key, "processed")
		if err != nil {
			log.Errorf("failed to set upload status: %v", fileErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
		}

		return &models.ResponseMsg{Message: "success"}, nil
	}
	log.Errorf("failed to find correct template")
	return nil, echo.NewHTTPError(http.StatusBadRequest, "неправильный шаблон документа, обратитесь к нам")
}

func (e ExcelServiceImpl) GetFileColumns(ctx context.Context, req *models.DirectusModel) ([]*models.FileColumns, error){
	if req.Collection != "uploads" {
		log.Warnf("collection is %s not uploads", req.Collection)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "collection is not uploads")
	}
	log.Infof("Start processing upload from directus: %v", req)
	// err := e.repo.SetUploadStatus(ctx, req.Key, "processing")
	// if err != nil {
	// 	log.Warnf("failed to set upload status", err)
	// 	return nil, err
	// }
	//time.Sleep(15 * time.Second) // todo удалить после демо 8.10
	_, uploadErr := e.repo.GetFromUploadCatalogue(ctx, req.Key)
	if uploadErr != nil {
		log.Warnf("failed to get upload catalog", uploadErr)
		return nil, uploadErr
	}

	endpoint := e.cfg.Aws.Host
	accessKeyID := e.cfg.Aws.SecretKey
	secretAccessKey := e.cfg.Aws.AccessKey
	bucket := e.cfg.Aws.Bucket
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(secretAccessKey, accessKeyID, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error("failed to connect to minio: ", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	minioObj, getObjErr := minioClient.GetObject(ctx, bucket, "upload.FileId", minio.GetObjectOptions{})
	if getObjErr != nil {
		log.Error("failed to get object:", getObjErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, getObjErr)
	}

	defer minioObj.Close()

	excelFile, fileErr := excelize.OpenReader(minioObj)
	if fileErr != nil {
		log.Errorf("failed to open reader: %v", fileErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fileErr)
	}

	rows, rowsErr := excelFile.GetRows(excelFile.GetSheetList()[0])

	if rowsErr != nil {
		log.Errorf("failed to read sheet: %v", rowsErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
	}
	
	var fileRows []*models.FileColumns
	for i, v := range rows[0]{
		fileRow := &models.FileColumns{}
		fileRow.RowId = int8(i)
		fileRow.RowName = v
		fileRows = append(fileRows, fileRow)
	}

	return fileRows, nil
}

// func (e ExcelServiceImpl) SaveCargoCatalogue(ctx context.Context, file *multipart.FileHeader) (*models.ResponseMsg, error) {
// 	src, err := file.Open()
// 	if err != nil {
// 		log.Errorf("failed ti open file: %v", err)
// 		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	defer src.Close()

// 	excelFile, fileErr := excelize.OpenReader(src)
// 	if fileErr != nil {
// 		log.Errorf("failed to open reader: %v", fileErr)
// 		return nil, echo.NewHTTPError(http.StatusBadRequest, fileErr)
// 	}
// 	fmt.Println(excelFile.GetSheetList())
// 	// Get all the rows in the Sheet1.
// 	fmt.Println("1")
// 	rows, rowsErr := excelFile.GetRows("Лист1")
// 	fmt.Println("1")
// 	if rowsErr != nil {
// 		log.Errorf("failed to read sheet: %v", rowsErr)
// 		return nil, echo.NewHTTPError(http.StatusBadRequest, "Не правильный наименование страницы excel файла. Переименуйте на Лист1")
// 	}

// 	for i, row := range rows {
// 		if i == 0 {
// 			continue
// 		}

// 		nomenclature := &models.Nomenclature{}
// 		nomenclature.CodeKsNsi = row[0]
// 		nomenclature.Link = row[1]
// 		nomenclature.CodeSkmtr = row[2]
// 		nomenclature.TmcMark = row[3]
// 		nomenclature.GostTu = row[4]
// 		if row[5] != "отсуствует" {
// 			nomenclature.DrawingName = row[5]
// 		}
// 		nomenclature.Measurement = row[6]
// 		nomenclature.CategoryName = row[7]
// 		nomenclature

// 		cargoCat := &models.CargoCatalogue{}
// 		cargoCat.MinLotShipment = row[8]

// 	}
// }
