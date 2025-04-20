package request

// CreatePaymentRequest represents a request to create a new payment
type CreatePaymentRequest struct {
	ContractorID string  `json:"contractor_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Description  string  `json:"description" binding:"required"`
}
