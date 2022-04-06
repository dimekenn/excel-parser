package models

type Category struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Type       string `json:"type"`
	ParentName string `json:"parent_name"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}

type Nomenclature struct {
	Id                    string                 `json:"id"`
	CodeSkmtr             string                 `json:"code_skmtr"`
	CodeKsNsi             string                 `json:"code_ks_nsi"`
	CodeAmto              string                 `json:"code_amto"`
	OKPD2                 string                 `json:"okpd2"`
	CodeTnved             string                 `json:"code_tnved"`
	Name                  string                 `json:"name"`
	TmcCodeVendor         string                 `json:"tmc_code_vendor"`
	TmcMark               string                 `json:"tmc_mark"`
	GostTu                string                 `json:"gost_tu"`
	DateOfManufacture     string                 `json:"date_of_manufacture"`
	Manufacturer          string                 `json:"manufacturer"`
	BatchNumber           string                 `json:"batch_number"`
	IsTax                 bool                   `json:"is_tax"`
	TaxPercentage         float32                `json:"tax_percentage"`
	PricePerUnit          float32                `json:"price_per_unit"`
	Measurement           string                 `json:"measurement"`
	PriceValidThrough     string                 `json:"price_valid_through"`
	WholesaleItems        *WholesaleItems        `json:"wholesale_items"`
	Quantity              int                    `json:"quantity"`
	ProductAvailability   bool                   `json:"product_availability"`
	HazardClass           string                 `json:"hazard_class"`
	PackagingType         string                 `json:"packaging_type"`
	PackingMaterial       string                 `json:"packing_material"`
	StorageType           string                 `json:"storage_type"`
	WeightNetto           float32                `json:"weight_netto"`
	WeightBrutto          float32                `json:"weight_brutto"`
	LoadingType           string                 `json:"loading_type"`
	WarehouseAddress      string                 `json:"warehouse_address"`
	Regions               string                 `json:"regions"`
	DeliveryType          string                 `json:"delivery_type"`
	PackageId             string                 `json:"package_id"`
	Length                float32                `json:"length"`
	Height                float32                `json:"height"`
	Width                 float32                `json:"width"`
	Volume                float32                `json:"volume"`
	AmountInPackage       int8                   `json:"amount_in_package"`
	Class                 string                 `json:"class"`
	Representation        string                 `json:"representation"`
	DeliveryAddress       string                 `json:"delivery_address"`
	FullName              string                 `json:"full_name"`
	Link                  string                 `json:"link"`
	Payload               *Mtr                   `json:"payload,omitempty"`
	DrawingName           string                 `json:"drawing_name"`
	CategoryName          string                 `json:"category_name"`
	OrganizerNomenclature *OrganizerNomenclature `json:"payload"`
	CompanyInn            string                 `json:"company_inn"`
	UserId                string                 `json:"user_id"`
	CargoCatalogue        *CargoCatalogue        `json:"cargo_catalogue"`
	PriceLists            []string               `json:"price"`
}

type Mtr struct {
	Link                     string `json:"link,omitempty"`
	DataVersion              string `json:"data_version,omitempty"`
	DeleteMark               string `json:"delete_mark,omitempty"`
	Code                     string `json:"code,omitempty"`
	Name                     string `json:"name,omitempty"`
	VendorCode               string `json:"vendor_code,omitempty"`
	Measurement              string `json:"measurement,omitempty"`
	Identifier               string `json:"identifier,omitempty"`
	CatalogueNumber          string `json:"catalogue_number,omitempty"`
	Class                    string `json:"class,omitempty"`
	Comments                 string `json:"comments,omitempty"`
	PropertySet              string `json:"property_set,omitempty"`
	TechDoc                  string `json:"tech_doc,omitempty"`
	Okved2                   string `json:"okved_2,omitempty"`
	Okpd2                    string `json:"okpd_2,omitempty"`
	Description              string `json:"description,omitempty"`
	FullName                 string `json:"full_name,omitempty"`
	SignOfUser               string `json:"sign_of_user,omitempty"`
	Manufacturer             string `json:"manufacturer,omitempty"`
	Tnved                    string `json:"tnved,omitempty"`
	DeleteRecord             string `json:"delete_record,omitempty"`
	DeleteItemType           string `json:"delete_item_type,omitempty"`
	DeleteRefPosition        string `json:"delete_ref_position,omitempty"`
	DeleteLayout             string `json:"delete_layout,omitempty"`
	SlAmto                   string `json:"sl_amto,omitempty"`
	SlManufacturerVendorCode string `json:"sl_manufacturer_vendor_code,omitempty"`
	SlManufacturerBarcode    string `json:"sl_manufacturer_barcode,omitempty"`
	SlDraw                   string `json:"sl_draw,omitempty"`
	SlWeightNetto            string `json:"sl_weight_netto,omitempty"`
	SlWeightBrutto           string `json:"sl_weight_brutto,omitempty"`
	SlPriority               string `json:"sl_priority,omitempty"`
	SlSupplierMeasurement    string `json:"sl_supplier_measurement,omitempty"`
	SlConversionFactor       string `json:"sl_conversion_factor,omitempty"`
	SlSupplierWeightNetto    string `json:"sl_supplier_weight_netto,omitempty"`
	SlSupplierWeightBrutto   string `json:"sl_supplier_weight_brutto,omitempty"`
	SlExpiryDate             string `json:"sl_expiry_date,omitempty"`
	SlManufacturerCountry    string `json:"sl_manufacturer_country,omitempty"`
	SlCheckInterval          string `json:"sl_check_interval,omitempty"`
	SlDrawingFile            string `json:"sl_drawing_file,omitempty"`
	SlImgFile                string `json:"sl_img_file,omitempty"`
	SlMarkTmc                string `json:"sl_mark_tmc,omitempty"`
	SlStateStandard          string `json:"sl_state_standard,omitempty"`
	SlPackage                string `json:"sl_package,omitempty"`
	SlHazardClass            string `json:"sl_hazard_class,omitempty"`
	SlNomenclatureSign       string `json:"sl_nomenclature_sign,omitempty"`
	SlSize                   string `json:"sl_size,omitempty"`
	MdmKey                   string `json:"mdm_key,omitempty"`
	NsiRequest               string `json:"nsi_request,omitempty"`
	NsiManualChange          string `json:"nsi_manual_change,omitempty"`
	Predefined               string `json:"predefined,omitempty"`
	PredefinedDataName       string `json:"predefined_data_name,omitempty"`
	Representation           string `json:"representation,omitempty"`
	Measurement1             string `json:"measurement_1,omitempty"`
	Coefficient              string `json:"coefficient,omitempty"`
	Purpose                  string `json:"purpose,omitempty"`
	Analog                   string `json:"analog,omitempty"`
	KindOfClassifier         string `json:"kind_of_classifier,omitempty"`
	Class1                   string `json:"class_1,omitempty"`
	Property                 string `json:"property,omitempty"`
	Value                    string `json:"value,omitempty"`
	TextString               string `json:"text_string,omitempty"`
	SparePart                string `json:"spare_part,omitempty"`
	Shipper                  string `json:"shipper,omitempty"`
	ShippingAddress          string `json:"shipping_address,omitempty"`
	MinShippingBatch         string `json:"min_shipping_batch,omitempty"`
	CharacteristicName       string `json:"characteristic_name,omitempty"`
	Characteristic           string `json:"characteristic,omitempty"`
	Value1                   string `json:"value_1,omitempty"`
}

type WholesaleItems struct {
	WholesalePricePerUnit float32 `json:"wholesale_price_per_unit"`
	WholesaleOrderFrom    int     `json:"wholesale_order_from"`
	WholesaleOrderTo      int     `json:"wholesale_order_to"`
}

type Company struct {
	Name   string `json:"name"`
	Inn    string `json:"inn"`
	UserId string `json:"user_id"`
}

type OrganizerNomenclature struct {
	NomenclatureCode                         string `json:"nomenclature_code"`
	NomenclatureType                         string `json:"nomenclature_type"`
	IsWeight                                 bool   `json:"is_weight"`
	WeightCoefficient                        string `json:"weight_coefficient"`
	WIPBalance                               string `json:"wip_balance"`
	PartitionAccountingBySeries              string `json:"partition_accounting_by_series"`
	AccountingBySeries                       string `json:"accounting_by_series"`
	KeepAccountingBySeriesWCD                string `json:"keep_accounting_by_series_wcd"`
	KeepAccountingAccordingToCharacteristics string `json:"keep_accounting_according_to_characteristics"`
	KindReproduction                         string `json:"kind_reproduction"`
	MainMeasurement                          string `json:"main_measurement"`
	ReportMeasurement                        string `json:"report_measurement"`
	ResidueMeasurement                       string `json:"residue_measurement"`
	Kit                                      string `json:"kit"`
	PurposeOfUse                             string `json:"purpose_of_use"`
	Comments                                 string `json:"comments"`
	Service                                  string `json:"service"`
	NomenclatureGroup                        string `json:"nomenclature_group"`
	FileImg                                  string `json:"file_img"`
	MainSupplier                             string `json:"main_supplier"`
	SalesManager                             string `json:"sales_manager"`
	ManufacturerCountry                      string `json:"manufacturer_country"`
	GTDNumber                                string `json:"gtd_number"`
	ArticleCost                              string `json:"article_cost"`
	RequiresExternalCertification            bool   `json:"requires_external_certification"`
	RequiresInternalCertification            bool   `json:"requires_internal_certification"`
	Set                                      bool   `json:"set"`
	OKP                                      string `json:"okp"`
	IsAlcohol                                bool   `json:"is_alcohol"`
	IsImportAlcohol                          bool   `json:"is_import_alcohol"`
	VolumeDAL                                string `json:"volume_dal"`
	QuarantineZone                           bool   `json:"quarantine_zone"`
	CodeSUMI                                 string `json:"code_sumi"`
	AMTOStatus                               string `json:"amto_status"`
	ENSKStatus                               string `json:"ensk_status"`
	ENSKName                                 string `json:"ensk_name"`
	ENSKTM                                   string `json:"ensktm"`
	ENSKBrandDesign                          string `json:"ensk_brand_design"`
	ENDSDefaultMark                          string `json:"ends_default_mark"`
	ENSKTechSpec                             string `json:"ensk_tech_spec"`
	ENSKMaterialMark                         string `json:"ensk_material_mark"`
	ENSKGostMaterial                         string `json:"ensk_gost_material"`
	CatalogueNumber                          string `json:"catalogue_number"`
	ENSKOKPClassificator                     string `json:"enskokp_classificator"`
	AMTONormName                             string `json:"amto_norm_name"`
	AMTOCodeForEOrder                        string `json:"amto_code_for_e_order"`
	ENSKExpertComments                       string `json:"ensk_expert_comments"`
	TMXClassificatorGP                       string `json:"tmx_classificator_gp"`
	TMXClassificatorOKP                      string `json:"tmx_classificator_okp"`
	TMXClassificatorRTK                      string `json:"tmx_classificator_rtk"`
	TMXCodePDM                               string `json:"tmx_code_pdm"`
	TMXItemType                              string `json:"tmx_item_type"`
	IsTobacco                                bool   `json:"is_tobacco"`
	IsShoes                                  bool   `json:"is_shoes"`
	TMXCodeMDM                               string `json:"tmx_code_mdm"`
}

type CargoCatalogue struct {
	MinLotShipment   string `json:"min_lot_shipment"`
	ManufacturerDays string `json:"manufacturer_days"`
}

// type GetExcelFromAwsByFileIdReq{
// 	FileId string `json:"file_id"`
// 	CompanyName string `json:"company_name"`
// }
