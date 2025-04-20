package request

// CreateTimelogRequest represents a request to create a new timelog entry
type CreateTimelogRequest struct {
	ContractorID string  `json:"contractor_id" binding:"required"`
	JobID        string  `json:"job_id" binding:"required"`
	Date         string  `json:"date" binding:"required"`
	Hours        float64 `json:"hours" binding:"required"`
	Description  string  `json:"description"`
}
