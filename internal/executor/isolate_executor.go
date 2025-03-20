package executor

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/utils"
	"go.uber.org/zap"
)

type isolateExecutor struct {
	baseExecutor
}

func (e *isolateExecutor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	submissionID := *submission.ID
	startTime := time.Now()

	e.logger.Info("Starting execution",
		zap.String("submission_id", submissionID),
		zap.Int("test_case_count", len(submission.TestCases)),
		zap.String("language", *submission.Language.SourceFileExt),
		zap.Int("time_limit_ms", submission.TimeLimitInMs),
		zap.Int("memory_limit_kb", submission.MemoryLimitInKb))

	command := e.buildCommand(submission)
	e.logger.Debug("Built execution command",
		zap.String("submission_id", submissionID),
		zap.String("command", command))

	defer func() {
		e.logger.Debug("Cleaning up files",
			zap.String("submission_id", submissionID),
			zap.String("isolate_dir", e.cfg.IsolateDir))
		e.cleanupFiles(submission)
	}()

	// Execute test cases in parallel
	var wg sync.WaitGroup
	wg.Add(len(submission.TestCases))

	e.logger.Debug("Launching test case executors",
		zap.String("submission_id", submissionID),
		zap.Int("parallelism", len(submission.TestCases)))

	for i, testcase := range submission.TestCases {
		go func(tc model.TestCase, idx int) {
			tcStartTime := time.Now()
			testCaseID := *tc.ID

			e.logger.Debug("Starting test case execution",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.Int("test_case_index", idx))

			executor := NewTestCaseExecutor(e.logger, e.cfg.IsolateDir, command, submission)
			result := executor.Execute(&tc)

			// Log test case result
			e.logger.Debug("Test case completed",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.String("status", *result.Status),
				zap.Float64("time_ms", result.TimeUsageInMs),
				zap.Float64("memory_kb", result.MemoryUsageInKb),
				zap.Duration("execution_duration", time.Since(tcStartTime)))

			ch <- result
			wg.Done()
		}(testcase, i)
	}

	wg.Wait()
	close(ch)

	e.logger.Info("Execution completed for all test cases",
		zap.String("submission_id", submissionID),
		zap.Duration("total_duration", time.Since(startTime)))

	return nil
}

func (e *isolateExecutor) Compile(submission *model.Submission) error {
	submissionID := *submission.ID
	startTime := time.Now()

	e.logger.Info("Starting compilation",
		zap.String("submission_id", submissionID),
		zap.String("language", *submission.Language.SourceFileExt))

	sourceFilename := e.writeSourceFile(submission)
	if sourceFilename == "" {
		e.logger.Error("Failed to write source file",
			zap.String("submission_id", submissionID),
			zap.String("isolate_dir", e.cfg.IsolateDir))
		return errors.New("failed to write source file")
	}

	e.logger.Debug("Source file written",
		zap.String("submission_id", submissionID),
		zap.String("filename", sourceFilename),
		zap.String("path", filepath.Join(e.cfg.IsolateDir, sourceFilename)))

	if !e.needsCompilation(submission) {
		e.logger.Info("No compilation needed for this language",
			zap.String("submission_id", submissionID),
			zap.String("language", *submission.Language.SourceFileExt))
		return nil
	}

	command := e.buildCompileCommand(submission, sourceFilename)
	e.logger.Debug("Built compile command",
		zap.String("submission_id", submissionID),
		zap.String("command", command))

	err := e.runCompileCommand(command, submissionID)
	if err != nil {
		e.logger.Error("Compilation failed",
			zap.String("submission_id", submissionID),
			zap.Error(err),
			zap.Duration("compile_time", time.Since(startTime)))
		return err
	}

	e.logger.Info("Compilation successful",
		zap.String("submission_id", submissionID),
		zap.Duration("compile_time", time.Since(startTime)))

	return nil
}

func (e *isolateExecutor) runCompileCommand(command string, submissionID string) error {
	startTime := time.Now()

	e.logger.Debug("Executing compile command",
		zap.String("submission_id", submissionID),
		zap.String("command", command),
		zap.String("working_dir", e.cfg.IsolateDir))

	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.cfg.IsolateDir
	var stderr strings.Builder
	execCmd.Stderr = &stderr

	if err := execCmd.Run(); err != nil {
		stderrOutput := stderr.String()
		truncatedErr := utils.TruncateString(stderrOutput, 500) // Limit very long error outputs

		e.logger.Error("Compilation process failed",
			zap.String("submission_id", submissionID),
			zap.Error(err),
			zap.String("stderr", truncatedErr),
			zap.Duration("duration", time.Since(startTime)))

		// Include full stderr output in the returned error
		return fmt.Errorf("%s", stderrOutput)
	}

	e.logger.Debug("Compilation process completed successfully",
		zap.String("submission_id", submissionID),
		zap.Duration("duration", time.Since(startTime)))

	return nil
}

func newIsolateExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	isolateExecutor := &isolateExecutor{
		baseExecutor: baseExecutor{
			logger: logger,
			cfg:    cfg,
		},
	}

	logger.Info("Created new isolate executor",
		zap.String("isolate_dir", cfg.IsolateDir))

	return isolateExecutor
}
