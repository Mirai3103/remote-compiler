package grpc

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/Mirai3103/remote-compiler/internal/executor"
	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/logger"
	"github.com/Mirai3103/remote-compiler/proto"
	"golang.org/x/sync/semaphore"
)

type ExecutionHandler struct {
	proto.UnimplementedExecutionServiceServer
	cfx              *config.Config
	compileSemaphore *semaphore.Weighted // Semaphore để giới hạn số lượng biên dịch
	executeSemaphore *semaphore.Weighted // Semaphore để giới hạn số lượng thực thi
	mu               sync.Mutex          // Mutex để bảo vệ truy cập vào semaphores
}

func NewExecutionHandler(cfx *config.Config) *ExecutionHandler {
	return &ExecutionHandler{
		cfx:              cfx,
		compileSemaphore: semaphore.NewWeighted(10), // Giới hạn 10 hoạt động biên dịch đồng thời
		executeSemaphore: semaphore.NewWeighted(15), // Giới hạn 15 hoạt động thực thi đồng thời
	}
}

func (h *ExecutionHandler) Execute(req *proto.Submission, stream proto.ExecutionService_ExecuteServer) error {
	log := logger.GetLogger()
	log.Info("Received a submission")

	submission := &model.Submission{
		ID:              &req.Id,
		Language:        convertLanguage(req.Language),
		Code:            &req.Code,
		TimeLimitInMs:   int(req.TimeLimitInMs),
		MemoryLimitInKb: int(req.MemoryLimitInKb),
		TestCases:       convertTestCases(req.TestCases),
		Settings:        convertSubmissionSettings(req.Settings),
	}

	var ex = executor.NewExecutor(log, h.cfx.Executor)

	// Sử dụng semaphore để giới hạn số lượng biên dịch đồng thời
	ctx := context.Background()
	if err := h.compileSemaphore.Acquire(ctx, 1); err != nil {
		log.Error("Failed to acquire compile semaphore: %v", zap.Error(err))
		return err
	}
	log.Info("Acquired compile semaphore, compiling submission")

	err := ex.Compile(submission)
	h.compileSemaphore.Release(1)
	log.Info("Released compile semaphore")

	if err == nil {
		// Biên dịch thành công, tiến hành thực thi
		ch := make(chan *model.SubmissionResult, len(submission.TestCases))

		// Acquire execute semaphore
		if err := h.executeSemaphore.Acquire(ctx, 1); err != nil {
			log.Error("Failed to acquire execute semaphore:", zap.Error(err))
			return err
		}
		log.Info("Acquired execute semaphore, executing submission")

		go func() {
			ex.Execute(submission, ch)
			h.executeSemaphore.Release(1)
			log.Info("Released execute semaphore")
		}()

		for result := range ch {
			stream.Send(&proto.SubmissionResult{
				SubmissionId:    *result.SubmissionID,
				TestCaseId:      *result.TestCaseID,
				Status:          *result.Status,
				Stdout:          *result.Stdout,
				MemoryUsageInKb: float32(result.MemoryUsageInKb),
				TimeUsageInMs:   float32(result.TimeUsageInMs),
			})
		}
	} else {
		// Biên dịch thất bại, gửi lỗi cho tất cả các test cases
		errStr := err.Error()
		for _, testCase := range submission.TestCases {
			_ = stream.Send(&proto.SubmissionResult{
				SubmissionId: *submission.ID,
				TestCaseId:   *testCase.ID,
				Status:       "Compile Error",
				Stdout:       errStr,
			})
		}
	}

	return nil
}

func convertLanguage(lang *proto.Language) *model.Language {
	if lang == nil {
		return nil
	}
	return &model.Language{
		SourceFileExt:  &lang.SourceFileExt,
		BinaryFileExt:  &lang.BinaryFileExt,
		CompileCommand: &lang.CompileCommand,
		RunCommand:     &lang.RunCommand,
	}
}

func convertTestCases(testCases []*proto.TestCase) []model.TestCase {
	var result []model.TestCase
	for _, testCase := range testCases {
		result = append(result, model.TestCase{
			ID:           &testCase.Id,
			Input:        &testCase.Input,
			ExpectOutput: &testCase.ExpectOutput,
			InputFile:    &testCase.InputFile,
			OutputFile:   &testCase.OutputFile,
		})
	}
	return result
}

func convertSubmissionSettings(settings *proto.SubmissionSettings) *model.SubmissionSettings {
	if settings == nil {
		return &model.SubmissionSettings{}
	}
	return &model.SubmissionSettings{
		WithTrim:          settings.WithTrim,
		WithCaseSensitive: settings.WithCaseSensitive,
		WithWhitespace:    settings.WithWhitespace,
	}
}
