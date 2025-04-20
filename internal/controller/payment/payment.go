package payment

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/internal/controller/payment/request"
	"github.com/mercor/payment-service/internal/payment/service"
)

type PaymentController struct {
	svc *service.PaymentService
}

var (
	ctrl     *PaymentController
	ctrlOnce sync.Once
)

func NewPaymentController(svc *service.PaymentService) *PaymentController {
	ctrlOnce.Do(func() {
		ctrl = &PaymentController{svc: svc}
	})
	return ctrl
}

// GET /api/v1/contractors/:contractor_id/payment-line-items?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (c *PaymentController) GetPaymentLineItemsForContractorPeriod(ctx *gin.Context) {
	contractorID := ctx.Param("contractor_id")
	startTime, err := strconv.ParseInt(ctx.Query("time_start"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	endTime, err := strconv.ParseInt(ctx.Query("time_end"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	items, err := c.svc.GetPaymentLineItemsForContractorPeriod(ctx, contractorID, startTime, endTime)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// PUT /api/v1/payment-line-items/:id
func (c *PaymentController) UpdatePaymentLineItemByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var updates *request.UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := c.svc.UpdatePaymentLineItemByID(ctx, id, updates); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}
