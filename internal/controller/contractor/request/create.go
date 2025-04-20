package request

// Request structs (similar to job/request)
type CreateContractorCtrlReq struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}
