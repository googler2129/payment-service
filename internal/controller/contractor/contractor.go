package contractor

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	svcreq "github.com/mercor/payment-service/internal/contractor/request"
	"github.com/mercor/payment-service/internal/controller/contractor/request"
	"github.com/mercor/payment-service/internal/domain"
)

type Controller struct {
	svc domain.ContractorServiceInterface
}

var (
	ctrl     *Controller
	ctrlOnce sync.Once
)

func NewController(svc domain.ContractorServiceInterface) *Controller {
	ctrlOnce.Do(func() {
		ctrl = &Controller{
			svc: svc,
		}
	})
	return ctrl
}

func (c *Controller) CreateContractor(ctx *gin.Context) {
	var req *request.CreateContractorCtrlReq

	// Binding and validation
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = c.svc.CreateContractor(ctx, convertCreateContractorCtrlReqToCreateContractorSvcReq(req))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func convertCreateContractorCtrlReqToCreateContractorSvcReq(req *request.CreateContractorCtrlReq) *svcreq.CreateContractorSvcReq {
	return &svcreq.CreateContractorSvcReq{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
}
