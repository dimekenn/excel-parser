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
}

type UploadsEntity struct {
	CompanyId string `json:"company_id"`
	FileId    string `json:"file_id"`
	UserId    string `json:"user_id"`
}
