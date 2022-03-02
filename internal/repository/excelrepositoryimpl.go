package repository

import (
	"context"
	"database/sql"
	"excel-service/internal/models"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	container "github.com/vielendanke/go-db-lb"
)

type ExcelRepositoryImpl struct {
	lb *container.LoadBalancer
}

func NewExcelRepository(lb *container.LoadBalancer) ExcelRepository {
	return &ExcelRepositoryImpl{lb: lb}
}

func (e ExcelRepositoryImpl) SaveNomenclature(ctx context.Context, nomenclature *models.Nomenclature, tx pgx.Tx) error {
	//var id string

	_, execErr := tx.Exec(
		ctx,
		//"insert into nomenclature (id, code_skmtr, code_ks_nsi, code_amto, okpd2, code_tnved, name, tmc_code_vendor, tmc_mark, date_of_manufacture, manufacturer, is_tax, tax_percentage, price_per_unit, measurement, price_valid_through, wholesale_price_per_unit, wholesale_order_from, wholesale_order_to, quantity, product_availability, hazard_class, packaging_type, packing_material, storage_type, weight_netto, weight_brutto, loading_type, warehouse_address, regions, delivery_type) values " +
		//	"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, (select id from measurement where name = $15), $16, $17, $18, $19, $20, $21, (select id from hazard_class where name = $22), (select id from packaging_type where name = $23), (select id from packing_material  where name = $24), (select id from storage_type where name = $25), $26, $27, (select id from loading_type  where name = $28), $29,(select id from regions where name = $30), (select id from delivery_type where name = $31)) returning id",
		"with nom as (insert into nomenclature (id, code_skmtr, code_ks_nsi, code_amto, okpd2, code_tnved, name, tmc_code_vendor, tmc_mark, gost_tu, date_of_manufacture, manufacturer, batch_number, is_tax, tax_percentage, price_per_unit, measurement, price_valid_through, wholesale_items, quantity, product_availability,  loading_type, regions, delivery_type) "+
			"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,  $11, $12, $13,  $14, $15, $16, (select id from measurement where name = $17), $18, $19, $20, $21, (select id from loading_type  where name = $22), (select id from regions where name = $23), (select id from delivery_type where name = $24)) returning id), "+
			"package as (insert into package(id, packaging_type, packing_material, name, storage_type, hazard_class, length, height, width, volume,  weight_brutto, weight_netto, amount_in_package) "+
			"values ($25, (select id from packaging_type where name = $26), (select id from packing_material  where name = $27), $28, (select id from storage_type  where name = $29), (select id from hazard_class where name = $30), $31, $32, $33, $34, $35, $36, $37) returning id) insert into nomenclature_package ( nomenclature_id, package_id) values ((select id from nom), (select id from package))",
		nomenclature.Id, newNullString(nomenclature.CodeSkmtr), newNullString(nomenclature.CodeKsNsi), newNullString(nomenclature.CodeAmto), nomenclature.OKPD2, nomenclature.CodeTnved, nomenclature.Name, newNullString(nomenclature.TmcCodeVendor), newNullString(nomenclature.TmcMark), newNullString(nomenclature.GostTu), newNullString(nomenclature.DateOfManufacture), newNullString(nomenclature.Manufacturer), newNullString(nomenclature.BatchNumber), nomenclature.IsTax, nomenclature.TaxPercentage, nomenclature.PricePerUnit, nomenclature.Measurement, nomenclature.PriceValidThrough, nomenclature.WholesaleItems, nomenclature.Quantity, nomenclature.ProductAvailability, nomenclature.LoadingType, nomenclature.Regions, nomenclature.DeliveryType, nomenclature.PackageId, nomenclature.PackagingType, nomenclature.PackingMaterial, nomenclature.Name, nomenclature.StorageType, nomenclature.HazardClass, nomenclature.Length, nomenclature.Height, nomenclature.Width, nomenclature.Volume, nomenclature.WeightBrutto, nomenclature.WeightNetto, nomenclature.AmountInPackage,
	)

	if execErr != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			log.Errorf("failed to roll back tx in SaveNomenclature: %v", rbErr)
			return echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("failed to insert nomenclature: %v", execErr)
		return echo.NewHTTPError(http.StatusInternalServerError, execErr)
	}

	return nil
}

func newNullString(s string) *sql.NullString {
	if len(s) == 0 {
		return &sql.NullString{}
	}
	return &sql.NullString{
		String: s,
		Valid:  true,
	}
}
