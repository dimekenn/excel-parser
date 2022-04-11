package models

type GetExcelFromAwsByFileIdReq struct {
	FileId      string `json:"file_id"`
	CompanyName string `json:"comany_name"`
}

type DirectusModel struct {
	Key        string `json:"key"`
	Collection string `json:"collection"`
	Accounting struct {
		User    string `json:"user"`
		Company string `json:"company"`
		Role    string `json:"role"`
	} `json:"accountability"`
	Payload struct {
		Files struct {
			DirectusFilesId []struct {
				Id string `json:"id"`
			} `json:"directus_files_id"`
		} `json:"files"`
		Catalogs struct {
			CatalogId []struct {
				Id string `json:"id"`
			} `json:"catalog_id"`
		} `json:"catalogs"`
	} `json:"payload"`
}

type UploadsEntity struct {
	CompanyId string `json:"company_id"`
	FileId    string `json:"file_id"`
	UserId    string `json:"user_id"`
}
