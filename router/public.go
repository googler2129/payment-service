package router

import (
	"context"

	"github.com/mercor/payment-service/internal/controller/contractor"
	"github.com/mercor/payment-service/internal/controller/job"
	"github.com/mercor/payment-service/internal/controller/payment"
	"github.com/mercor/payment-service/internal/controller/timelog"
	"github.com/mercor/payment-service/pkg/cluster"
	uhttp "github.com/mercor/payment-service/pkg/http"
)

func PublicRoutes(ctx context.Context, s *uhttp.Server) (err error) {
	paymentController, _ := payment.Wire(ctx, cluster.GetCluster().DbCluster)
	timelogController, _ := timelog.Wire(ctx, cluster.GetCluster().DbCluster)
	contractorController, _ := contractor.Wire(ctx, cluster.GetCluster().DbCluster)
	jobController, _ := job.Wire(ctx, cluster.GetCluster().DbCluster)

	contractor := s.Engine.Group("/api/v1/contractors")
	{
		contractor.POST("", contractorController.CreateContractor)
	}

	job := s.Engine.Group("/api/v1/jobs")
	{
		job.POST("", jobController.CreateJob)
		job.GET("/extended", jobController.GetJobsByStatus)
		job.GET("/active/:contractor_id", jobController.GetActiveJobsForContractor)
	}

	paymentLineItems := s.Engine.Group("/api/v1/payment-line-items")
	{
		paymentLineItems.PUT(":id", paymentController.UpdatePaymentLineItemByID)
	}

	payment := s.Engine.Group("/api/v1/contractors/:contractor_id/payment-line-items")
	{
		payment.GET("", paymentController.GetPaymentLineItemsForContractorPeriod)
	}

	timelog := s.Engine.Group("/api/v1/contractors/:contractor_id/timelogs")
	{
		timelog.GET("", timelogController.GetTimelogsForContractorPeriod)
	}

	return nil
}
