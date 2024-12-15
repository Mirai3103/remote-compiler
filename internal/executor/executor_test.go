package executor_test

import (
	"testing"

	"github.com/Mirai3103/remote-compiler/internal/executor"
	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func ptr(s string) *string {
	return &s
}

func TestExecutor(t *testing.T) {
	cfg, err := config.LoadConfig("/root/remote-compiler/config.yaml")
	assert.NoError(t, err)

	log := logger.GetLogger()
	executor := executor.NewExecutor(log, cfg.Executor)

	submission := &model.Submission{
		Language: &model.Language{
			Version:        ptr("g++ 12.2.0"),
			Name:           ptr("C++"),
			SourceFileExt:  ptr(".cpp"),
			BinaryFileExt:  ptr(".out"),
			CompileCommand: ptr("g++ $SourceFileName -o $BinaryFileName"),
			RunCommand:     ptr("$BinaryFileName"),
		},
		ID:          ptr("1"),
		Code:        ptr("#include <iostream>\nint main() { std::cout << \"Hello, World!\"; return 0; }"),
		TimeLimit:   1,
		MemoryLimit: 32000,
		TestCases: []model.TestCase{
			{
				ID:           ptr("1"),
				Input:        ptr(""),
				ExpectOutput: ptr("Hello, World!"),
			},
		},
	}

	err = executor.Compile(submission)
	assert.NoError(t, err)

	ch := make(chan *model.SubmissionResult, len(submission.TestCases))
	err = executor.Execute(submission, ch)
	assert.NoError(t, err)
	var result *model.SubmissionResult
	for r := range ch {
		result = r
		log.Sugar().Info(zap.Any("result", result))
	}
}
