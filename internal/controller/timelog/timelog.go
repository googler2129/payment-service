package timelog

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/internal/timelog/service"
)

type TimelogController struct {
	svc *service.TimelogService
}

var (
	ctrl     *TimelogController
	ctrlOnce sync.Once
)

func NewTimelogController(svc *service.TimelogService) *TimelogController {
	ctrlOnce.Do(func() {
		ctrl = &TimelogController{svc: svc}
	})
	return ctrl
}

// GET /api/v1/contractors/:contractor_id/timelogs?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (c *TimelogController) GetTimelogsForContractorPeriod(ctx *gin.Context) {
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

	timelogs, err := c.svc.GetTimelogsForContractorPeriod(ctx, contractorID, startTime, endTime)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, timelogs)
}
