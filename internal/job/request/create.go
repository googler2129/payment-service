package request

type CreateJobSvcReq struct {
	Status       string  `json:"status"`
	Rate         float64 `json:"rate"`
	Title        string  `json:"title"`
	CompanyID    string  `json:"company_id"`
	ContractorID string  `json:"contractor_id"`
}
