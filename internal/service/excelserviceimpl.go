package service

import (
	"context"
	"excel-service/internal/configs"
	"excel-service/internal/models"
	"excel-service/internal/repository"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	container "github.com/vielendanke/go-db-lb"
	"github.com/xuri/excelize/v2"
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
		taxPercentage, taxErr := strconv.ParseFloat(row[14], 8)
		if taxErr != nil {
			log.Errorf("failed to parse string to float: %v", taxErr)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат процента налога")
		}
		nomenclature.TaxPercentage = float32(taxPercentage)
		pricePerUnit, unitPriceErr := strconv.ParseFloat(row[15], 8)
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

		//wNetto, wNettoErr := strconv.Atoi(row[37])
		//if wNettoErr != nil {
		//	log.Errorf("failed to parse string to int: %v", wNettoErr)
		//	return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(нетто)")
		//}
		//wBrutto, wBruttoErr := strconv.Atoi(row[38])
		//if wBruttoErr != nil {
		//	log.Errorf("failed to parse string to int: %v", wBruttoErr)
		//	return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(брутто)")
		//}
		//nomenclature.WeightNetto = int16(wNetto)
		//nomenclature.WeightBrutto = int16(wBrutto)
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

	// tx, txErr := e.lb.CallPrimaryPreferred().PGxPool().Begin(ctx)
	// if txErr != nil {
	// 	log.Errorf("failed to begin tx: %v", txErr)
	// 	return nil, echo.NewHTTPError(http.StatusInternalServerError, txErr)
	// }
	// defer func(ctx context.Context) {
	// 	cErr := tx.Commit(ctx)
	// 	if cErr != nil {
	// 		log.Errorf("failed to commit tx in service: %v", cErr)
	// 		return
	// 	}
	// }(ctx)

	//var nomenclatures []*models.Nomenclature

	for i := 195841; i <= len(rows); i++ {
		v := rows[i]
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
				return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(нетто): "+v[29])
			}
			nomenclature.WeightNetto = float32(wNetto)
		}

		if v[30] != "" {
			wBrutto, wBruttoErr := strconv.ParseFloat(v[30], 8)
			if wBruttoErr != nil {
				log.Errorf("failed to parse string to int: %v", wBruttoErr)
				return nil, echo.NewHTTPError(http.StatusBadRequest, "не правильный формат вес(брутто)"+v[30])
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

		err := e.repo.SaveNomenclature(ctx, nomenclature, nil)
		if err != nil {
			return nil, err
		}

		//nomenclatures = append(nomenclatures, nomenclature)
	}

	//for _, v := range nomenclatures {
	//	err := e.repo.SaveNomenclature(ctx, v, tx)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return &models.ResponseMsg{Message: "success"}, nil
}

func takeVolume(name string) (string, float32, float32, float32) {
	regV, regErr := regexp.Compile(`\d+х.+х\d+`)
	if regErr != nil {
		log.Errorf("failed to regexp compile: %v", regErr)
		return name, 0, 0, 0
	}
	volume := regV.FindString(name)
	if volume != "" {
		vols := strings.Split(volume, "х")
		if len(vols) == 3 {
			length, lenErr := strconv.ParseFloat(strings.Replace(vols[0], ",", ".", 1), 32)
			if lenErr != nil {
				log.Errorf("failed to parse float: %v", lenErr)
				return name, 0, 0, 0
			}
			height, heiErr := strconv.ParseFloat(strings.Replace(vols[1], ",", ".", 1), 32)
			if heiErr != nil {
				log.Errorf("failed to parse float: %v", heiErr)
				return name, 0, 0, 0
			}
			width, widErr := strconv.ParseFloat(strings.Replace(vols[2], ",", ".", 1), 32)
			if widErr != nil {
				log.Errorf("failed to parse float: %v", widErr)
				return name, 0, 0, 0
			}
			name = strings.Replace(name, volume, "", 1)
			return name, float32(length), float32(height), float32(width)
		} else {
			return name, 0, 0, 0
		}
	}
	return name, 0, 0, 0
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

	//var nomenclatures []*models.Nomenclature

	for i := 590297; i < len(rows); i++ {
		// if i < 1 {
		// 	continue
		// }
		row := rows[i]
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
			nomenclature.UserId = "Supplier"
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
		err := e.repo.SaveNomenclature(ctx, nomenclature, tx)
		if err != nil {
			log.Error(err)
			e.repo.NewErrorNomenclatureId(ctx, i)
		}
	}

	// saveErr := e.repo.SaveArrayNomenclature(ctx, nomenclatures, tx)
	// if saveErr != nil {
	// 	fmt.Println("save array nom errors: ", saveErr)
	// 	return nil, saveErr
	// }
	return &models.ResponseMsg{Message: "success"}, nil
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
	_, err = os.Create("asd/json")
	if err != nil {
		log.Error("create err", err)
		return nil, err
	}
	log.Info("prowel")
	err = minioClient.FGetObject(ctx, bucket, req.FileId, req.FileId, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("get object err:", err)
		return nil, err
	}

	return &models.ResponseMsg{Message: "success"}, nil
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
