package executor

import (
	"errors"
	"os/exec"
	"strings"
	"sync"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"go.uber.org/zap"
)

type isolateExecutor struct {
	baseExecutor
}

func (e *isolateExecutor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	command := e.buildCommand(submission)
	defer e.cleanupFiles(submission)

	var wg sync.WaitGroup
	for _, testcase := range submission.TestCases {
		wg.Add(1)
		go func(tc model.TestCase) {
			defer wg.Done()
			executor := NewTestCaseExecutor(e.logger, e.cfg.IsolateDir, command, submission)
			ch <- executor.Execute(&tc)
		}(testcase)
	}

	wg.Wait()
	close(ch)
	return nil
}

func (e *isolateExecutor) Compile(submission *model.Submission) error {
	sourceFilename := e.writeSourceFile(submission)
	if sourceFilename == "" {
		return errors.New("failed to write source file")
	}
	if !e.needsCompilation(submission) {
		return nil
	}

	command := e.buildCompileCommand(submission, sourceFilename)
	return e.runCompileCommand(command)
}

func (e *isolateExecutor) runCompileCommand(command string) error {
	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.cfg.IsolateDir
	var stderr strings.Builder
	execCmd.Stderr = &stderr

	if err := execCmd.Run(); err != nil {
		e.logger.Error("Compilation error",
			zap.Error(err),
			zap.String("stderr", stderr.String()))
		return errors.New(stderr.String())
	}
	return nil
}

func newIsolateExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	return &isolateExecutor{
		baseExecutor: baseExecutor{
			logger: logger,
			cfg:    cfg,
		},
	}
}
