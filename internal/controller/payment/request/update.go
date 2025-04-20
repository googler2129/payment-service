package request

type UpdatePaymentRequest struct {
	JobUID     string  `json:"job_uid" binding:"required"`
	TimelogUID string  `json:"timelog_uid" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Status     string  `json:"status" binding:"required"`
}
