package models

type ResponseMsg struct {
	Message string `json:"message"`
}

type Nomenclature struct {
	Id                  string          `json:"id"`
	CodeSkmtr           string          `json:"code_skmtr"`
	CodeKsNsi           string          `json:"code_ks_nsi"`
	CodeAmto            string          `json:"code_amto"`
	OKPD2               string          `json:"okpd2"`
	CodeTnved           string          `json:"code_tnved"`
	Name                string          `json:"name"`
	TmcCodeVendor       string          `json:"tmc_code_vendor"`
	TmcMark             string          `json:"tmc_mark"`
	GostTu              string          `json:"gost_tu"`
	DateOfManufacture   string          `json:"date_of_manufacture"`
	Manufacturer        string          `json:"manufacturer"`
	BatchNumber         string          `json:"batch_number"`
	IsTax               bool            `json:"is_tax"`
	TaxPercentage       float32         `json:"tax_percentage"`
	PricePerUnit        float32         `json:"price_per_unit"`
	Measurement         string          `json:"measurement"`
	PriceValidThrough   string          `json:"price_valid_through"`
	WholesaleItems      *WholesaleItems `json:"wholesale_items"`
	Quantity            int             `json:"quantity"`
	ProductAvailability bool            `json:"product_availability"`
	HazardClass         string          `json:"hazard_class"`
	PackagingType       string          `json:"packaging_type"`
	PackingMaterial     string          `json:"packing_material"`
	StorageType         string          `json:"storage_type"`
	WeightNetto         int16           `json:"weight_netto"`
	WeightBrutto        int16           `json:"weight_brutto"`
	LoadingType         string          `json:"loading_type"`
	WarehouseAddress    string          `json:"warehouse_address"`
	Regions             string          `json:"regions"`
	DeliveryType        string          `json:"delivery_type"`
	PackageId           string          `json:"package_id"`
	Length              float32         `json:"length"`
	Height              float32         `json:"height"`
	Width               float32         `json:"width"`
	Volume              float32         `json:"volume"`
	AmountInPackage     int8            `json:"amount_in_package"`
}

type WholesaleItems struct {
	WholesalePricePerUnit float32 `json:"wholesale_price_per_unit"`
	WholesaleOrderFrom    int     `json:"wholesale_order_from"`
	WholesaleOrderTo      int     `json:"wholesale_order_to"`
}
