package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/utils"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

// riskExecutor runs code directly on the host machine without sandboxing (not recommended for production)
type riskExecutor struct {
	baseExecutor
}

// Compile implements Executor.
func (r *riskExecutor) Compile(submission *model.Submission) error {
	submissionID := *submission.ID
	startTime := time.Now()

	r.logger.Info("Starting compilation in risk executor",
		zap.String("submission_id", submissionID),
		zap.String("language", *submission.Language.SourceFileExt),
		zap.String("os", runtime.GOOS))

	sourceFilename := r.writeSourceFile(submission)
	if sourceFilename == "" {
		r.logger.Error("Failed to write source file",
			zap.String("submission_id", submissionID),
			zap.String("base_dir", r.cfg.IsolateDir))
		return errors.New("failed to write source file")
	}

	r.logger.Debug("Source file written",
		zap.String("submission_id", submissionID),
		zap.String("filename", sourceFilename))

	if !r.needsCompilation(submission) {
		r.logger.Info("No compilation needed for this language",
			zap.String("submission_id", submissionID),
			zap.String("language", *submission.Language.SourceFileExt))
		return nil
	}

	command := r.buildCompileCommand(submission, sourceFilename)
	r.logger.Debug("Built compile command",
		zap.String("submission_id", submissionID),
		zap.String("command", command))

	err := r.runCompileCommand(command, submissionID)

	if err != nil {
		r.logger.Error("Compilation failed",
			zap.String("submission_id", submissionID),
			zap.Error(err),
			zap.Duration("duration", time.Since(startTime)))
		return err
	}

	r.logger.Info("Compilation successful",
		zap.String("submission_id", submissionID),
		zap.Duration("compile_time", time.Since(startTime)))

	return nil
}

func (r *riskExecutor) runCompileCommand(command string, submissionID string) error {
	startTime := time.Now()
	osName := runtime.GOOS

	r.logger.Debug("Executing compile command",
		zap.String("submission_id", submissionID),
		zap.String("os", osName),
		zap.String("command", command))

	var cmd *exec.Cmd
	if osName == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		truncatedOutput := utils.TruncateString(outputStr, 500)

		r.logger.Error("Compilation process failed",
			zap.String("submission_id", submissionID),
			zap.Error(err),
			zap.String("output", truncatedOutput),
			zap.Duration("duration", time.Since(startTime)))

		return fmt.Errorf("compilation error: %s\n%s", err.Error(), outputStr)
	}

	r.logger.Debug("Compilation process completed",
		zap.String("submission_id", submissionID),
		zap.Duration("duration", time.Since(startTime)),
		zap.String("output_size", fmt.Sprintf("%d bytes", len(output))))

	return nil
}

// Execute implements Executor.
func (r *riskExecutor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	submissionID := *submission.ID
	startTime := time.Now()
	osName := runtime.GOOS
	maxTimeInMs := submission.TimeLimitInMs
	maxMemoryInBytes := submission.MemoryLimitInKb * 1024

	r.logger.Info("Starting execution in risk executor",
		zap.String("submission_id", submissionID),
		zap.String("os", osName),
		zap.Int("test_case_count", len(submission.TestCases)),
		zap.Int("time_limit_ms", maxTimeInMs),
		zap.Int("memory_limit_kb", submission.MemoryLimitInKb))

	command := r.buildCommand(submission)
	r.logger.Debug("Built execution command",
		zap.String("submission_id", submissionID),
		zap.String("command", command))

	defer func() {
		r.logger.Debug("Cleaning up files",
			zap.String("submission_id", submissionID))
		r.cleanupFiles(submission)
	}()

	var wg sync.WaitGroup
	var succeededTests int32
	var failedTests int32

	for i, testcase := range submission.TestCases {
		wg.Add(1)
		go func(tc model.TestCase, index int) {
			defer wg.Done()
			tcStartTime := time.Now()
			testCaseID := *tc.ID

			r.logger.Debug("Starting test case execution",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.Int("test_case_index", index))

			result := &model.SubmissionResult{
				SubmissionID:    submission.ID,
				TestCaseID:      tc.ID,
				MemoryUsageInKb: 0,
				TimeUsageInMs:   0,
			}

			var cmd *exec.Cmd
			if osName == "windows" {
				cmd = exec.Command("cmd", "/C", command)
			} else {
				cmd = exec.Command("bash", "-c", command)
			}
			cmd.Env = os.Environ()

			// Setup pipes for input and output
			stdin, err := cmd.StdinPipe()
			if err != nil {
				r.logger.Error("Failed to get stdin pipe",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Error(err))

				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			// Create buffers for output
			var outBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = &outBuf

			// Start the command
			err = cmd.Start()
			if err != nil {
				r.logger.Error("Failed to start command",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Error(err))

				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			pid := int32(cmd.Process.Pid)
			r.logger.Debug("Process started",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.Int32("pid", pid))

			// Channel to signal process completion
			done := make(chan error, 1)

			// Memory monitoring goroutine
			memoryMonitorCtx, memoryMonitorCancel := context.WithCancel(context.Background())
			defer memoryMonitorCancel()

			var maxMemUsage uint64 = 0
			memoryExceeded := false

			go func() {
				ticker := time.NewTicker(10 * time.Millisecond)
				defer ticker.Stop()

				for {
					select {
					case <-memoryMonitorCtx.Done():
						return
					case <-ticker.C:
						proc, err := process.NewProcess(pid)
						if err != nil {
							continue // Process might have exited
						}

						memInfo, err := proc.MemoryInfo()
						if err != nil {
							continue
						}

						if memInfo != nil {
							currentMemory := memInfo.RSS
							if currentMemory > maxMemUsage {
								maxMemUsage = currentMemory
								r.logger.Debug("Memory usage updated",
									zap.String("submission_id", submissionID),
									zap.String("test_case_id", testCaseID),
									zap.Uint64("memory_bytes", currentMemory),
									zap.Float64("memory_mb", float64(currentMemory)/(1024*1024)))
							}

							// Check if memory limit exceeded
							if uint64(maxMemoryInBytes) > 0 && currentMemory > uint64(maxMemoryInBytes) {
								memoryExceeded = true
								r.logger.Info("Memory limit exceeded",
									zap.String("submission_id", submissionID),
									zap.String("test_case_id", testCaseID),
									zap.Uint64("usage_bytes", currentMemory),
									zap.Float64("usage_mb", float64(currentMemory)/(1024*1024)),
									zap.Int("limit_kb", submission.MemoryLimitInKb))

								// Kill process
								if cmd.Process != nil {
									cmd.Process.Kill()
								}
								return
							}
						}
					}
				}
			}()

			// Write input and close stdin
			r.logger.Debug("Writing input to process",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.Int("input_length", len(*tc.Input)))

			_, err = stdin.Write([]byte(*tc.Input))
			if err != nil {
				r.logger.Error("Failed to write to stdin",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Error(err))

				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			stdin.Close()

			// Setup timeout context
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxTimeInMs)*time.Millisecond)
			defer cancel()

			execStartTime := time.Now()

			// Wait for process in a goroutine
			go func() {
				done <- cmd.Wait()
			}()

			// Wait for process completion, timeout, or memory limit
			var waitErr error
			select {
			case waitErr = <-done:
				// Process completed normally
				r.logger.Debug("Process completed",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Duration("runtime", time.Since(execStartTime)))
			case <-ctx.Done():
				// Time limit exceeded
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				r.logger.Info("Time limit exceeded",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Int("time_limit_ms", maxTimeInMs),
					zap.Duration("runtime", time.Since(execStartTime)))

				timeExceededMsg := "time limit exceeded"
				result.Stdout = &timeExceededMsg
				result.Status = &StatusTimeLimitExceeded
				result.TimeUsageInMs = float64(maxTimeInMs)
				result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0 // Convert to KB
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			// Stop memory monitoring
			memoryMonitorCancel()

			// Check if memory limit was exceeded
			if memoryExceeded {
				r.logger.Info("Result: Memory limit exceeded",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Float64("memory_used_kb", float64(maxMemUsage)/1024.0))

				memExceededMsg := "memory limit exceeded"
				result.Stdout = &memExceededMsg
				result.Status = &StatusMemoryLimitExceeded
				result.TimeUsageInMs = float64(time.Since(execStartTime).Microseconds()) / 1000.0 // ms
				result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0                            // Convert to KB
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			duration := time.Since(execStartTime)
			result.TimeUsageInMs = float64(duration.Microseconds()) / 1000.0 // ms
			result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0           // Convert to KB

			if waitErr != nil {
				outputStr := outBuf.String()
				truncatedOutput := utils.TruncateString(outputStr, 200)

				r.logger.Error("Command execution failed",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.String("output", truncatedOutput),
					zap.Error(waitErr))

				errStr := fmt.Sprintf("%s\n%s", waitErr.Error(), outputStr)
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				atomic.AddInt32(&failedTests, 1)
				ch <- result
				return
			}

			// Process output
			outputStr := outBuf.String()
			result.Stdout = &outputStr

			if compare(outputStr, *tc.ExpectOutput, submission.Settings) {
				result.Status = &StatusSuccess
				r.logger.Info("Test case passed",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Float64("time_ms", result.TimeUsageInMs),
					zap.Float64("memory_kb", result.MemoryUsageInKb))
				atomic.AddInt32(&succeededTests, 1)
			} else {
				result.Status = &StatusWrongAnswer
				r.logger.Info("Test case failed: Wrong answer",
					zap.String("submission_id", submissionID),
					zap.String("test_case_id", testCaseID),
					zap.Float64("time_ms", result.TimeUsageInMs),
					zap.Float64("memory_kb", result.MemoryUsageInKb))
				atomic.AddInt32(&failedTests, 1)

				if r.logger.Core().Enabled(zap.DebugLevel) {
					r.logger.Debug("Output comparison",
						zap.String("submission_id", submissionID),
						zap.String("test_case_id", testCaseID),
						zap.String("expected", utils.TruncateString(*tc.ExpectOutput, 100)),
						zap.String("actual", utils.TruncateString(outputStr, 100)))
				}
			}

			ch <- result
			r.logger.Debug("Test case execution completed",
				zap.String("submission_id", submissionID),
				zap.String("test_case_id", testCaseID),
				zap.Duration("duration", time.Since(tcStartTime)))
		}(testcase, i)
	}

	wg.Wait()
	close(ch)

	r.logger.Info("Execution completed for all test cases",
		zap.String("submission_id", submissionID),
		zap.Int32("succeeded", succeededTests),
		zap.Int32("failed", failedTests),
		zap.Duration("total_duration", time.Since(startTime)))

	return nil
}

func newRiskExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	logger.Info("Creating risk executor",
		zap.String("base_dir", cfg.IsolateDir),
		zap.String("os", runtime.GOOS),
		zap.String("warning", "This executor runs code directly on the host machine without sandboxing"))

	return &riskExecutor{
		baseExecutor: baseExecutor{
			logger: logger,
			cfg:    cfg,
		},
	}
}
