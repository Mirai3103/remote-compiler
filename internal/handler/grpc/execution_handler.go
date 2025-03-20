package grpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Mirai3103/remote-compiler/internal/executor"
	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/logger"
	"github.com/Mirai3103/remote-compiler/pkg/utils"
	"github.com/Mirai3103/remote-compiler/proto"
	"golang.org/x/sync/semaphore"
)

type ExecutionHandler struct {
	proto.UnimplementedExecutionServiceServer
	cfx              *config.Config
	compileSemaphore *semaphore.Weighted
	executeSemaphore *semaphore.Weighted
}

func NewExecutionHandler(cfx *config.Config) *ExecutionHandler {
	handler := &ExecutionHandler{
		cfx:              cfx,
		compileSemaphore: semaphore.NewWeighted(int64(cfx.Executor.MaxCompileConcurrent)),
		executeSemaphore: semaphore.NewWeighted(int64(cfx.Executor.MaxExecuteConcurrent)),
	}

	log := logger.GetLogger()
	log.Info("Initialized execution handler",
		zap.Int("max_compile_concurrent", cfx.Executor.MaxCompileConcurrent),
		zap.Int("max_execute_concurrent", cfx.Executor.MaxExecuteConcurrent))

	return handler
}

func (h *ExecutionHandler) Execute(req *proto.Submission, stream proto.ExecutionService_ExecuteServer) error {
	log := logger.GetLogger()
	reqID := req.Id
	startTime := time.Now()

	log.Info("Execution request received",
		zap.String("submission_id", reqID),
		zap.String("language", req.Language.SourceFileExt),
		zap.Int32("time_limit_ms", req.TimeLimitInMs),
		zap.Int32("memory_limit_kb", req.MemoryLimitInKb),
		zap.Int("test_cases_count", len(req.TestCases)))

	if log.Core().Enabled(zap.DebugLevel) {
		log.Debug("Submission details",
			zap.String("submission_id", reqID),
			zap.String("code_preview", utils.TruncateString(req.Code, 100)),
			zap.Any("settings", req.Settings))
	}

	submission := &model.Submission{
		ID:              &req.Id,
		Language:        convertLanguage(req.Language),
		Code:            &req.Code,
		TimeLimitInMs:   int(req.TimeLimitInMs),
		MemoryLimitInKb: int(req.MemoryLimitInKb),
		TestCases:       convertTestCases(req.TestCases),
		Settings:        convertSubmissionSettings(req.Settings),
	}

	// Create execution context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var ex = executor.NewExecutor(log, h.cfx.Executor)
	log.Debug("Executor initialized", zap.String("submission_id", reqID))

	// Acquire compile semaphore
	log.Debug("Waiting for compile semaphore", zap.String("submission_id", reqID))
	semaphoreAcquireStart := time.Now()

	if err := h.compileSemaphore.Acquire(ctx, 1); err != nil {
		log.Error("Failed to acquire compile semaphore",
			zap.Error(err),
			zap.String("submission_id", reqID),
			zap.Duration("wait_time", time.Since(semaphoreAcquireStart)))
		return fmt.Errorf("semaphore acquisition failed: %w", err)
	}

	log.Info("Compiling submission",
		zap.String("submission_id", reqID),
		zap.Duration("semaphore_wait", time.Since(semaphoreAcquireStart)))

	// Time the compilation
	compileStart := time.Now()
	err := ex.Compile(submission)
	compileDuration := time.Since(compileStart)

	// Release semaphore
	h.compileSemaphore.Release(1)
	log.Debug("Released compile semaphore", zap.String("submission_id", reqID))

	if err != nil {
		log.Error("Compilation failed",
			zap.Error(err),
			zap.String("submission_id", reqID),
			zap.Duration("compile_time", compileDuration))

		// Handle compilation error by sending error to all test cases
		errStr := err.Error()
		for _, testCase := range submission.TestCases {
			_ = stream.Send(&proto.SubmissionResult{
				SubmissionId: *submission.ID,
				TestCaseId:   *testCase.ID,
				Status:       "Compile Error",
				Stdout:       errStr,
			})
			log.Debug("Sent compile error result",
				zap.String("submission_id", reqID),
				zap.String("test_case_id", *testCase.ID))
		}

		log.Info("Execution completed with compile error",
			zap.String("submission_id", reqID),
			zap.Duration("total_time", time.Since(startTime)))
		return nil
	}

	log.Info("Compilation successful",
		zap.String("submission_id", reqID),
		zap.Duration("compile_time", compileDuration))

	// Execution phase
	log.Debug("Waiting for execute semaphore", zap.String("submission_id", reqID))
	execSemaphoreStart := time.Now()

	if err := h.executeSemaphore.Acquire(ctx, 1); err != nil {
		log.Error("Failed to acquire execute semaphore",
			zap.Error(err),
			zap.String("submission_id", reqID),
			zap.Duration("wait_time", time.Since(execSemaphoreStart)))
		return fmt.Errorf("execute semaphore acquisition failed: %w", err)
	}

	log.Info("Executing submission",
		zap.String("submission_id", reqID),
		zap.Duration("semaphore_wait", time.Since(execSemaphoreStart)),
		zap.Int("test_cases", len(submission.TestCases)))

	// Execute test cases
	execStart := time.Now()
	ch := make(chan *model.SubmissionResult, len(submission.TestCases))

	go func() {
		ex.Execute(submission, ch)
		execDuration := time.Since(execStart)
		h.executeSemaphore.Release(1)
		log.Info("Execution completed",
			zap.String("submission_id", reqID),
			zap.Duration("execution_time", execDuration))
	}()

	// Process results
	resultCount := 0
	sucessCount := 0

	for result := range ch {
		resultCount++

		isSuccess := *result.Status == "Accepted"
		if isSuccess {
			sucessCount++
		}

		log.Debug("Test case result",
			zap.String("submission_id", *result.SubmissionID),
			zap.String("test_case_id", *result.TestCaseID),
			zap.String("status", *result.Status),
			zap.Float64("time_ms", result.TimeUsageInMs),
			zap.Float64("memory_kb", result.MemoryUsageInKb))

		stream.Send(&proto.SubmissionResult{
			SubmissionId:    *result.SubmissionID,
			TestCaseId:      *result.TestCaseID,
			Status:          *result.Status,
			Stdout:          *result.Stdout,
			MemoryUsageInKb: float32(result.MemoryUsageInKb),
			TimeUsageInMs:   float32(result.TimeUsageInMs),
		})
	}

	log.Info("All results sent",
		zap.String("submission_id", reqID),
		zap.Int("total_test_cases", resultCount),
		zap.Int("successful_test_cases", sucessCount),
		zap.Duration("total_processing_time", time.Since(startTime)))

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
