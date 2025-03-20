package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing/iotest"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/isolate"

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

	boxId := boxIDManager.Acquire()
	defer boxIDManager.Release(boxId)
	if boxId == -1 {
		return e.handleError(fmt.Errorf("no available box id"), boxId, testcase.ID)
	}
	e.boxId = boxId
	initBoxMutex.Lock()

	if err := e.setupEnvironment(testcase, boxId); err != nil {
		initBoxMutex.Unlock()
		e.logger.Error("Init box failed",
			zap.Int("boxId", boxId),
			zap.Error(err))

		return e.handleError(err, boxId, testcase.ID)
	}
	initBoxMutex.Unlock()
	defer e.cleanup(boxId, testcase)

	execResult, err := e.runIsolatedCommand(testcase, boxId)
	if err != nil {
		e.logger.Error("Run isolated command failed",
			zap.Int("boxId", boxId),
			zap.Error(err))
		return e.handleError(err, boxId, testcase.ID)
	}

	return e.processExecutionResult(execResult, testcase, e.submission.Settings)
}

func (e *TestCaseExecutor) setupEnvironment(testcase *model.TestCase, boxId int) error {
	inputFilename := e.isolateDir + "/" + testcase.GetInputFileName()
	commandShFile := e.isolateDir + "/" + snowflakeid.NewString() + ".sh"
	e.shFile = &commandShFile

	if err := os.WriteFile(inputFilename, []byte(*testcase.Input), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(commandShFile, []byte(e.command), 0644); err != nil {
		return err
	}

	_, err := isolate.InitBox(boxId)
	return err
}

func (e *TestCaseExecutor) cleanup(boxId int, testcase *model.TestCase) {
	isolate.CleanBox(boxId)
	os.Remove(e.isolateDir + "/" + testcase.GetInputFileName())
	os.Remove(e.isolateDir + "/" + testcase.GetExpectOutputFileName())
	os.Remove(e.isolateDir + "/" + testcase.GetExpectOutputFileName() + ".meta")
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
	e.logger.Info("Run isolated command", zap.String("command", strCmd))
	execCmd := exec.Command(args[0], args[1:]...)
	execCmd.Dir = e.isolateDir
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout
	if err := execCmd.Run(); err != nil {
		output, err1 := e.readOutput(testcase.GetExpectOutputFileName())
		if err1 == nil {
			e.logger.Error("Run isolated command failed with output", zap.String("output", output))
			return nil, fmt.Errorf("%s", output)
		}
		var metaResult *isolate.MetaResult
		if metaResult, err1 = isolate.NewMetaResultFromFile(metaOutFilename); err1 == nil {
			e.logger.Error("Run isolated command failed with meta", zap.Any("metaResult", metaResult))
		}

		return nil, err
	}

	return isolate.NewMetaResultFromFile(metaOutFilename)
}

func (e *TestCaseExecutor) processExecutionResult(metaResult *isolate.MetaResult, testcase *model.TestCase, settings *model.SubmissionSettings) *model.SubmissionResult {
	result := &model.SubmissionResult{
		SubmissionID:    e.submission.ID,
		TestCaseID:      testcase.ID,
		MemoryUsageInKb: metaResult.CGMem,
		TimeUsageInMs:   metaResult.TimeWall,
	}

	if metaResult.ExitCode != 0 {
		result.Status = &StatusRuntimeError
		result.Stdout = &metaResult.Message
		return result
	}

	if metaResult.TimeWall > float64(e.submission.TimeLimitInMs) {
		result.Status = &StatusTimeLimitExceeded
		return result
	}

	if metaResult.CGMem > float64(e.submission.MemoryLimitInKb) {
		result.Status = &StatusMemoryLimitExceeded
		return result
	}

	output, err := e.readOutput(testcase.GetExpectOutputFileName())
	if err != nil {
		result.Status = &StatusRuntimeError
		errStr := err.Error()
		result.Stdout = &errStr
		return result
	}

	e.logger.Info("compare output", zap.String("output", output), zap.String("expectOutput", *testcase.ExpectOutput))
	if !compare(output, *testcase.ExpectOutput, settings) {
		result.Status = &StatusWrongAnswer
		result.Stdout = &output
		return result
	}

	result.Status = &StatusSuccess
	result.Stdout = &output
	return result
}

func (e *TestCaseExecutor) readOutput(filename string) (string, error) {
	file, err := os.Open("/var/local/lib/isolate/" + fmt.Sprint(e.boxId) + "/box/" + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := io.ReadAll(iotest.OneByteReader(file))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *TestCaseExecutor) handleError(err error, boxId int, testcaseId *string) *model.SubmissionResult {

	errStr := err.Error()
	return &model.SubmissionResult{
		SubmissionID: e.submission.ID,
		Status:       &StatusRuntimeError,
		Stdout:       &errStr,
		TestCaseID:   testcaseId,
	}
}
