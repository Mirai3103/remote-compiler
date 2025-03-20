package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing/iotest"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/isolate"
	"github.com/Mirai3103/remote-compiler/pkg/utils"

	snowflakeid "github.com/Mirai3103/remote-compiler/pkg/snowflake_id"
	"go.uber.org/zap"
)

var initBoxMutex = &sync.Mutex{}

type TestCaseExecutor struct {
	logger                *zap.Logger
	isolateDir            string
	isolateCommandBuilder *isolate.IsolateCommandBuilder
	command               string
	submission            *model.Submission
	shFile                *string
	boxId                 int
}

func NewTestCaseExecutor(logger *zap.Logger, isolateDir string, command string, submission *model.Submission) *TestCaseExecutor {
	logger = logger.With(
		zap.String("component", "TestCaseExecutor"),
		zap.String("submissionId", *submission.ID),
		zap.String("language", *submission.Language.SourceFileExt),
	)

	logger.Info("Creating new test case executor",
		zap.Int("timeLimit", submission.TimeLimitInMs),
		zap.Int("memoryLimit", submission.MemoryLimitInKb),
	)

	return &TestCaseExecutor{
		logger:     logger,
		isolateDir: isolateDir,
		command:    command,
		submission: submission,

		isolateCommandBuilder: isolate.NewIsolateCommandBuilder().
			WithProcesses(4).
			WithWallTime((submission.TimeLimitInMs / 1000) + 4).
			WithMaxFileSize(5120).
			AddDir("/etc:noexec").
			AddDir(isolateDir).
			WithCGroup().
			WithTime(submission.TimeLimitInMs / 1000).
			WithExtraTime(submission.TimeLimitInMs / 1000).
			WithCGroupMemory(submission.MemoryLimitInKb).
			WithStackSize(submission.MemoryLimitInKb).
			AddEnv("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").
			WithStderrToStdout(),
	}
}

func (e *TestCaseExecutor) Execute(testcase *model.TestCase) *model.SubmissionResult {
	executionStart := time.Now()
	e.logger.Info("Starting test case execution",
		zap.String("testCaseId", *testcase.ID),
		zap.Int("inputSize", len(*testcase.Input)),
		zap.Int("expectedOutputSize", len(*testcase.ExpectOutput)),
	)

	boxId := boxIDManager.Acquire()
	if boxId == -1 {
		e.logger.Error("Failed to acquire box id",
			zap.String("testCaseId", *testcase.ID))
		return e.handleError(fmt.Errorf("no available box id"), boxId, testcase.ID)
	}

	e.logger.Debug("Acquired box id",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	e.boxId = boxId
	defer boxIDManager.Release(boxId)

	initBoxMutex.Lock()
	e.logger.Debug("Acquired init box mutex",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	if err := e.setupEnvironment(testcase, boxId); err != nil {
		initBoxMutex.Unlock()
		e.logger.Error("Environment setup failed",
			zap.Int("boxId", boxId),
			zap.String("testCaseId", *testcase.ID),
			zap.Error(err),
			zap.String("isolateDir", e.isolateDir))

		return e.handleError(err, boxId, testcase.ID)
	}

	e.logger.Debug("Environment setup successful",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID),
		zap.String("inputFile", e.isolateDir+"/"+testcase.GetInputFileName()),
		zap.String("commandFile", *e.shFile))

	initBoxMutex.Unlock()
	e.logger.Debug("Released init box mutex",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	defer func() {
		cleanupStart := time.Now()
		e.cleanup(boxId, testcase)
		e.logger.Debug("Cleanup complete",
			zap.Int("boxId", boxId),
			zap.String("testCaseId", *testcase.ID),
			zap.Duration("cleanupDuration", time.Since(cleanupStart)))
	}()

	execStart := time.Now()
	e.logger.Info("Running isolated command",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	execResult, err := e.runIsolatedCommand(testcase, boxId)
	if err != nil {
		e.logger.Error("Isolated command execution failed",
			zap.Int("boxId", boxId),
			zap.String("testCaseId", *testcase.ID),
			zap.Error(err),
			zap.Duration("executionDuration", time.Since(execStart)))
		return e.handleError(err, boxId, testcase.ID)
	}

	e.logger.Info("Isolated command execution successful",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID),
		zap.Duration("executionDuration", time.Since(execStart)),
		zap.Float64("timeUsage", execResult.TimeWall),
		zap.Float64("memoryUsage", execResult.CGMem),
		zap.Int("exitCode", execResult.ExitCode),
		zap.String("status", execResult.Status))

	result := e.processExecutionResult(execResult, testcase, e.submission.Settings)

	e.logger.Info("Test case execution completed",
		zap.String("testCaseId", *testcase.ID),
		zap.String("status", *result.Status),
		zap.Float64("timeUsageMs", result.TimeUsageInMs),
		zap.Float64("memoryUsageKb", result.MemoryUsageInKb),
		zap.Duration("totalExecutionTime", time.Since(executionStart)))

	return result
}

func (e *TestCaseExecutor) setupEnvironment(testcase *model.TestCase, boxId int) error {
	setupStart := time.Now()
	e.logger.Debug("Setting up environment",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	inputFilename := e.isolateDir + "/" + testcase.GetInputFileName()
	commandShFile := e.isolateDir + "/" + snowflakeid.NewString() + ".sh"
	e.shFile = &commandShFile

	if err := os.WriteFile(inputFilename, []byte(*testcase.Input), 0644); err != nil {
		e.logger.Error("Failed to write input file",
			zap.String("filename", inputFilename),
			zap.Error(err))
		return err
	}

	e.logger.Debug("Input file written successfully",
		zap.String("filename", inputFilename),
		zap.Int("size", len(*testcase.Input)))

	if err := os.WriteFile(commandShFile, []byte(e.command), 0644); err != nil {
		e.logger.Error("Failed to write command file",
			zap.String("filename", commandShFile),
			zap.Error(err))
		return err
	}

	e.logger.Debug("Command file written successfully",
		zap.String("filename", commandShFile),
		zap.Int("size", len(e.command)))

	_, err := isolate.InitBox(boxId)
	if err != nil {
		e.logger.Error("Failed to initialize isolate box",
			zap.Int("boxId", boxId),
			zap.Error(err))
		return err
	}

	e.logger.Debug("Environment setup completed",
		zap.Int("boxId", boxId),
		zap.Duration("setupDuration", time.Since(setupStart)))

	return nil
}

func (e *TestCaseExecutor) cleanup(boxId int, testcase *model.TestCase) {
	e.logger.Debug("Starting cleanup",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))

	if err := isolate.CleanBox(boxId); err != nil {
		e.logger.Warn("Failed to clean isolate box",
			zap.Int("boxId", boxId),
			zap.Error(err))
	}

	filesToRemove := []string{
		e.isolateDir + "/" + testcase.GetInputFileName(),
		e.isolateDir + "/" + testcase.GetExpectOutputFileName(),
		e.isolateDir + "/" + testcase.GetExpectOutputFileName() + ".meta",
	}

	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			e.logger.Warn("Failed to remove file during cleanup",
				zap.String("filename", file),
				zap.Error(err))
		} else {
			e.logger.Debug("Removed file during cleanup",
				zap.String("filename", file))
		}
	}

	e.logger.Debug("Cleanup completed",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcase.ID))
}

func (e *TestCaseExecutor) runIsolatedCommand(testcase *model.TestCase, boxId int) (*isolate.MetaResult, error) {
	outputFilename := testcase.GetExpectOutputFileName()
	metaOutFilename := e.isolateDir + "/" + testcase.GetExpectOutputFileName() + ".meta"
	inputFilename := e.isolateDir + "/" + testcase.GetInputFileName()
	commandShFile := *e.shFile

	args := e.isolateCommandBuilder.Clone().
		WithBoxID(boxId).
		WithStdinFile(inputFilename).
		WithStdoutFile(outputFilename).
		WithMetaFile(metaOutFilename).
		WithRunCommands("/bin/bash", commandShFile).
		Build()
	strCmd := strings.Join(args, " ")

	e.logger.Debug("Built isolate command",
		zap.String("command", strCmd),
		zap.String("metaFile", metaOutFilename),
		zap.String("inputFile", inputFilename),
		zap.String("outputFile", outputFilename))

	execCmd := exec.Command(args[0], args[1:]...)
	execCmd.Dir = e.isolateDir

	// Capture stderr and stdout for logging
	var stdoutBuf, stderrBuf strings.Builder
	execCmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	execCmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

	startTime := time.Now()
	if err := execCmd.Run(); err != nil {
		execDuration := time.Since(startTime)
		e.logger.Error("Isolated command execution failed",
			zap.String("command", strCmd),
			zap.Error(err),
			zap.Duration("executionDuration", execDuration),
			zap.String("stdout", stdoutBuf.String()),
			zap.String("stderr", stderrBuf.String()))

		output, err1 := e.readOutput(testcase.GetExpectOutputFileName())
		if err1 == nil {
			e.logger.Error("Command failed with output in box",
				zap.String("output", output))
			return nil, fmt.Errorf("%s", output)
		}

		var metaResult *isolate.MetaResult
		if metaResult, err1 = isolate.NewMetaResultFromFile(metaOutFilename); err1 == nil {
			e.logger.Error("Command failed but meta file exists",
				zap.Any("metaResult", metaResult))
		} else {
			e.logger.Error("Failed to read meta file after command failure",
				zap.String("metaFile", metaOutFilename),
				zap.Error(err1))
		}

		return nil, err
	}

	execDuration := time.Since(startTime)
	e.logger.Debug("Isolated command execution successful",
		zap.String("command", strCmd),
		zap.Duration("executionDuration", execDuration),
		zap.String("stdout", stdoutBuf.String()),
		zap.String("stderr", stderrBuf.String()))

	metaResult, err := isolate.NewMetaResultFromFile(metaOutFilename)
	if err != nil {
		e.logger.Error("Failed to read meta result file",
			zap.String("metaFile", metaOutFilename),
			zap.Error(err))
		return nil, err
	}

	e.logger.Debug("Read meta result file successfully",
		zap.String("metaFile", metaOutFilename),
		zap.Any("metaResult", metaResult))

	return metaResult, nil
}

func (e *TestCaseExecutor) processExecutionResult(metaResult *isolate.MetaResult, testcase *model.TestCase, settings *model.SubmissionSettings) *model.SubmissionResult {
	e.logger.Debug("Processing execution result",
		zap.String("testCaseId", *testcase.ID),
		zap.Float64("timeWall", metaResult.TimeWall),
		zap.Float64("cgMem", metaResult.CGMem),
		zap.Int("exitCode", metaResult.ExitCode),
		zap.String("status", metaResult.Status))

	result := &model.SubmissionResult{
		SubmissionID:    e.submission.ID,
		TestCaseID:      testcase.ID,
		MemoryUsageInKb: metaResult.CGMem,
		TimeUsageInMs:   metaResult.TimeWall,
	}

	if metaResult.ExitCode != 0 {
		e.logger.Info("Submission result: Runtime Error",
			zap.String("testCaseId", *testcase.ID),
			zap.Int("exitCode", metaResult.ExitCode),
			zap.String("message", metaResult.Message))

		result.Status = &StatusRuntimeError
		result.Stdout = &metaResult.Message
		return result
	}

	if metaResult.TimeWall > float64(e.submission.TimeLimitInMs) {
		e.logger.Info("Submission result: Time Limit Exceeded",
			zap.String("testCaseId", *testcase.ID),
			zap.Float64("timeUsed", metaResult.TimeWall),
			zap.Int("timeLimit", e.submission.TimeLimitInMs))

		result.Status = &StatusTimeLimitExceeded
		return result
	}

	if metaResult.CGMem > float64(e.submission.MemoryLimitInKb) {
		e.logger.Info("Submission result: Memory Limit Exceeded",
			zap.String("testCaseId", *testcase.ID),
			zap.Float64("memoryUsed", metaResult.CGMem),
			zap.Int("memoryLimit", e.submission.MemoryLimitInKb))

		result.Status = &StatusMemoryLimitExceeded
		return result
	}

	output, err := e.readOutput(testcase.GetExpectOutputFileName())
	if err != nil {
		e.logger.Error("Failed to read output file",
			zap.String("testCaseId", *testcase.ID),
			zap.String("filename", testcase.GetExpectOutputFileName()),
			zap.Error(err))

		result.Status = &StatusRuntimeError
		errStr := err.Error()
		result.Stdout = &errStr
		return result
	}

	e.logger.Debug("Comparing output with expected output",
		zap.String("testCaseId", *testcase.ID),
		zap.String("output", utils.TruncateString(output, 500)),
		zap.String("expectOutput", utils.TruncateString(*testcase.ExpectOutput, 500)),
		zap.Int("outputLength", len(output)),
		zap.Int("expectOutputLength", len(*testcase.ExpectOutput)))

	if !compare(output, *testcase.ExpectOutput, settings) {
		e.logger.Info("Submission result: Wrong Answer",
			zap.String("testCaseId", *testcase.ID),
			zap.Int("outputLength", len(output)),
			zap.Int("expectOutputLength", len(*testcase.ExpectOutput)))

		result.Status = &StatusWrongAnswer
		result.Stdout = &output
		return result
	}

	e.logger.Info("Submission result: Success",
		zap.String("testCaseId", *testcase.ID),
		zap.Float64("timeUsed", metaResult.TimeWall),
		zap.Float64("memoryUsed", metaResult.CGMem))

	result.Status = &StatusSuccess
	result.Stdout = &output
	return result
}

func (e *TestCaseExecutor) readOutput(filename string) (string, error) {
	boxPath := "/var/local/lib/isolate/" + fmt.Sprint(e.boxId) + "/box/" + filename
	e.logger.Debug("Reading output file", zap.String("path", boxPath))

	file, err := os.Open(boxPath)
	if err != nil {
		e.logger.Error("Failed to open output file",
			zap.String("path", boxPath),
			zap.Error(err))
		return "", err
	}
	defer file.Close()

	b, err := io.ReadAll(iotest.OneByteReader(file))
	if err != nil {
		e.logger.Error("Failed to read output file",
			zap.String("path", boxPath),
			zap.Error(err))
		return "", err
	}

	e.logger.Debug("Successfully read output file",
		zap.String("path", boxPath),
		zap.Int("bytes", len(b)))

	return string(b), nil
}

func (e *TestCaseExecutor) handleError(err error, boxId int, testcaseId *string) *model.SubmissionResult {
	e.logger.Error("Handling execution error",
		zap.Int("boxId", boxId),
		zap.String("testCaseId", *testcaseId),
		zap.Error(err))

	errStr := err.Error()
	return &model.SubmissionResult{
		SubmissionID: e.submission.ID,
		Status:       &StatusRuntimeError,
		Stdout:       &errStr,
		TestCaseID:   testcaseId,
	}
}
