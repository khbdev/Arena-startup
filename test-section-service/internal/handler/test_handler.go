package handler

import (
	"context"
	"test-section-serve/internal/service"

	"github.com/khbdev/arena-startup-proto/proto/test-section"
)

type ResultHandler struct {
	testService *service.TestService
	test_sectionpb.UnimplementedResultServiceServer
}

func NewResultHandler(testService *service.TestService) *ResultHandler {
	return &ResultHandler{
		testService: testService,
	}
}

func (h *ResultHandler) GetUserTestResult(ctx context.Context, req *test_sectionpb.GetUserTestResultRequest) (*test_sectionpb.GetUserTestResultResponse, error) {

	telegramID := req.GetTelegramId()
	testID := req.GetTestId()

	
	resultJSON, err := h.testService.GetUserTestResult(telegramID, testID)
	if err != nil {
		return nil, err
	}


	return &test_sectionpb.GetUserTestResultResponse{
		JsonData: resultJSON,
	}, nil
}
