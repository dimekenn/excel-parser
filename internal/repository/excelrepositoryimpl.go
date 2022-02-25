package repository

import (
	"context"
	"database/sql"
	"excel-service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	container "github.com/vielendanke/go-db-lb"
	"net/http"
)

type ExcelRepositoryImpl struct {
	lb *container.LoadBalancer
}

func NewExcelRepository(lb *container.LoadBalancer) ExcelRepository {
	return &ExcelRepositoryImpl{lb: lb}
}

func (e ExcelRepositoryImpl) SaveNomenclature(ctx context.Context, nomenclature *models.Nomenclature, tx pgx.Tx) (string, error) {
	var id string

	execErr := tx.QueryRow(
		ctx,
		"insert into nomenclature (id, code_skmtr, code_ks_nsi, code_amto, okpd2, code_tnved, name, tmc_code_vendor, tmc_mark, date_of_manufacture, manufacturer, is_tax, tax_percentage, price_per_unit, measurement, price_valid_through, wholesale_price_per_unit, wholesale_order_from, wholesale_order_to, quantity, product_availability, hazard_class, packaging_type, packing_material, storage_type, weight_netto, weight_brutto, loading_type, warehouse_address, regions, delivery_type) values " +
			"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, (select id from measurement where name = $15), $16, $17, $18, $19, $20, $21, (select id from hazard_class where name = $22), (select id from packaging_type where name = $23), (select id from packing_material  where name = $24), (select id from storage_type where name = $25), $26, $27, (select id from loading_type  where name = $28), $29,(select id from regions where name = $30), (select id from delivery_type where name = $31)) returning id",
		//query,
		nomenclature.Id, newNullString(nomenclature.CodeSkmtr), newNullString(nomenclature.CodeKsNsi), newNullString(nomenclature.CodeAmto), nomenclature.OKPD2, nomenclature.CodeTnved, nomenclature.Name, newNullString(nomenclature.TmcCodeVendor), newNullString(nomenclature.TmcMark), newNullString(nomenclature.DateOfManufacture), newNullString(nomenclature.Manufacturer), nomenclature.IsTax, nomenclature.TaxPercentage, nomenclature.PricePerUnit, nomenclature.Measurement, nomenclature.PriceValidThrough, nomenclature.WholesalePricePerUnit, nomenclature.WholesaleOrderFrom, nomenclature.WholesaleOrderTo, nomenclature.Quantity, nomenclature.ProductAvailability, nomenclature.HazardClass, nomenclature.PackagingType, nomenclature.PackingMaterial, nomenclature.StorageType, nomenclature.WeightNetto, nomenclature.WeightBrutto, nomenclature.LoadingType, nomenclature.WarehouseAddress, nomenclature.Regions, nomenclature.DeliveryType,
	).Scan(&id)

	if execErr != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			log.Errorf("failed to roll back tx in SaveNomenclature: %v", rbErr)
			return "", echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("failed to insert nomenclature: %v", execErr)
		return "", echo.NewHTTPError(http.StatusInternalServerError, execErr)
	}

	return id, nil
}

func newNullString(s string) *sql.NullString {
	if len(s) == 0 {
		return &sql.NullString{}
	}
	return &sql.NullString{
		String: s,
		Valid: true,
	}
}
