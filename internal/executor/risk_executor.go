package executor

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

// chạy trực tiếp trên máy chủ, không cần tạo sandbox (không khuyến khích)
type riskExecutor struct {
	baseExecutor
}

// Compile implements Executor.
func (r *riskExecutor) Compile(submission *model.Submission) error {
	sourceFilename := r.writeSourceFile(submission)
	if sourceFilename == "" {
		return errors.New("failed to write source file")
	}
	if !r.needsCompilation(submission) {
		return nil
	}
	command := r.buildCompileCommand(submission, sourceFilename)
	return r.runCompileCommand(command)
}

func (r *riskExecutor) runCompileCommand(command string) error {

	osName := runtime.GOOS
	var cmd *exec.Cmd
	if osName == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.logger.Error("failed to compile", zap.String("output", string(output)), zap.Error(err))
		return err
	}
	return nil

}

// Execute implements Executor.
// Execute implements Executor.
func (r *riskExecutor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	defer r.cleanupFiles(submission)
	command := r.buildCommand(submission)
	var wg sync.WaitGroup
	osName := runtime.GOOS
	maxTimeInMs := submission.TimeLimitInMs
	maxMemoryInBytes := submission.MemoryLimitInKb * 1024

	for _, testcase := range submission.TestCases {
		wg.Add(1)
		go func(tc model.TestCase) {
			defer wg.Done()
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
				r.logger.Error("failed to get stdin pipe", zap.Error(err))
				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
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
				r.logger.Error("failed to start command", zap.Error(err))
				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				ch <- result
				return
			}

			pid := int32(cmd.Process.Pid)

			// Create context with timeout

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
							}

							// Check if memory limit exceeded
							if uint64(maxMemoryInBytes) > 0 && currentMemory > uint64(maxMemoryInBytes) {
								memoryExceeded = true
								r.logger.Info("memory limit exceeded",
									zap.Uint64("usage", currentMemory),
									zap.Int("limit", maxMemoryInBytes))

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
			_, err = stdin.Write([]byte(*tc.Input))

			if err != nil {
				r.logger.Error("failed to write to stdin", zap.Error(err))
				errStr := err.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				ch <- result
				return
			}

			stdin.Close()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxTimeInMs)*time.Millisecond)
			defer cancel()
			startTime := time.Now()
			// Wait for process in a goroutine
			go func() {
				done <- cmd.Wait()
			}()

			// Wait for process completion, timeout, or memory limit
			var waitErr error
			select {
			case waitErr = <-done:
				// Process completed normally
			case <-ctx.Done():
				// Time limit exceeded
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				timeExceededMsg := "time limit exceeded"
				result.Stdout = &timeExceededMsg
				result.Status = &StatusTimeLimitExceeded
				result.TimeUsageInMs = float64(maxTimeInMs)
				result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0 // Convert to KB
				ch <- result
				return
			}

			// Stop memory monitoring
			memoryMonitorCancel()

			// Check if memory limit was exceeded
			if memoryExceeded {
				memExceededMsg := "memory limit exceeded"
				result.Stdout = &memExceededMsg
				result.Status = &StatusMemoryLimitExceeded
				result.TimeUsageInMs = float64(time.Since(startTime).Microseconds()) / 1000.0 // ms
				result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0                        // Convert to KB
				ch <- result
				return
			}

			duration := time.Since(startTime)
			result.TimeUsageInMs = float64(duration.Microseconds()) / 1000.0 // ms
			result.MemoryUsageInKb = float64(maxMemUsage) / 1024.0           // Convert to KB

			if waitErr != nil {
				r.logger.Error("command execution failed", zap.String("output", outBuf.String()), zap.Error(waitErr))
				errStr := waitErr.Error()
				result.Stdout = &errStr
				result.Status = &StatusRuntimeError
				ch <- result
				return
			}

			// Process output
			outputStr := outBuf.String()
			result.Stdout = &outputStr
			if compare(outputStr, *testcase.ExpectOutput, submission.Settings) {
				result.Status = &StatusSuccess
			} else {
				result.Status = &StatusWrongAnswer
			}
			ch <- result
		}(testcase)
	}

	wg.Wait()
	close(ch)
	return nil
}

func newRiskExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	return &riskExecutor{
		baseExecutor: baseExecutor{
			logger: logger,
			cfg:    cfg,
		},
	}
}
