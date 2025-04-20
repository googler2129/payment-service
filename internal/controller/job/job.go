package job

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/internal/controller/job/request"
	"github.com/mercor/payment-service/internal/domain"
	svcreq "github.com/mercor/payment-service/internal/job/request"
)

type Controller struct {
	svc domain.JobServiceInterface
}

var (
	ctrl     *Controller
	ctrlOnce sync.Once
)

func NewController(svc domain.JobServiceInterface) *Controller {
	ctrlOnce.Do(func() {
		ctrl = &Controller{
			svc: svc,
		}
	})
	return ctrl
}

func (c *Controller) CreateJob(ctx *gin.Context) {
	var req *request.CreateJobCtrlReq

	// Binding and validation
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = c.svc.CreateJob(ctx, convertCreateJobCtrlReqToCreateJobSvcReq(req))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) GetJobsByStatus(ctx *gin.Context) {
	jobs, err := c.svc.GetJobsByStatus(ctx, "extended")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, jobs)
}

func (c *Controller) GetActiveJobsForContractor(ctx *gin.Context) {
	contractorID := ctx.Param("contractor_id")
	jobs, err := c.svc.GetActiveJobsForContractor(ctx, contractorID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, jobs)
}

func convertCreateJobCtrlReqToCreateJobSvcReq(req *request.CreateJobCtrlReq) *svcreq.CreateJobSvcReq {
	return &svcreq.CreateJobSvcReq{
		Status:       req.Status,
		Rate:         req.Rate,
		Title:        req.Title,
		CompanyID:    req.CompanyID,
		ContractorID: req.ContractorID,
	}
}
