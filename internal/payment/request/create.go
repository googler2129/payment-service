package request

// Add payment service request structs here as needed.

// CreatePaymentRequest represents a request to create a new payment
type CreatePaymentRequest struct {
	ContractorID string  `json:"contractor_id"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Date         string  `json:"date"`
}
