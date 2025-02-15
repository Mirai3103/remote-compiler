package grpc

import (
	"github.com/Mirai3103/remote-compiler/internal/executor"
	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/logger"
	"github.com/Mirai3103/remote-compiler/proto"
)

type ExecutionHandler struct {
	proto.UnimplementedExecutionServiceServer
	cfx *config.Config
}

func NewExecutionHandler(cfx *config.Config) *ExecutionHandler {
	return &ExecutionHandler{
		cfx: cfx,
	}
}

func (h *ExecutionHandler) Execute(req *proto.Submission, stream proto.ExecutionService_ExecuteServer) error {
	log := logger.GetLogger()
	log.Info("Received a submission")
	submission := &model.Submission{
		ID:          &req.Id,
		Language:    convertLanguage(req.Language),
		Code:        &req.Code,
		TimeLimit:   int(req.TimeLimit),
		MemoryLimit: int(req.MemoryLimit),
		TestCases:   convertTestCases(req.TestCases),
		Settings:    convertSubmissionSettings(req.Settings),
	}
	var ex = executor.NewExecutor(log, h.cfx.Executor)
	err := ex.Compile(submission)
	if err == nil {
		ch := make(chan *model.SubmissionResult, len(submission.TestCases))
		go ex.Execute(submission, ch)
		for result := range ch {
			stream.Send(&proto.SubmissionResult{
				SubmissionId: *result.SubmissionID,
				TestCaseId:   *result.TestCaseID,
				Status:       *result.Status,
				Stdout:       *result.Stdout,
				MemoryUsage:  float32(result.MemoryUsage),
				TimeUsage:    float32(result.TimeUsage),
			})
		}
	} else {
		errStr := err.Error()
		for _, testCase := range submission.TestCases {
			stream.Send(&proto.SubmissionResult{
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
	return &model.SubmissionSettings{
		WithTrim:          settings.WithTrim,
		WithCaseSensitive: settings.WithCaseSensitive,
		WithWhitespace:    settings.WithWhitespace,
	}
}
