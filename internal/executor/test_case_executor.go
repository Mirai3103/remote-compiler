package executor

import (
	"io"
	"os"
	"os/exec"
	"testing/iotest"

	"github.com/Mirai3103/remote-compiler/internal/isolate"
	model "github.com/Mirai3103/remote-compiler/internal/models"
	snowflakeid "github.com/Mirai3103/remote-compiler/pkg/snowflake_id"
	"go.uber.org/zap"
)

type TestCaseExecutor struct {
	logger                *zap.Logger
	isolateDir            string
	isolateCommandBuilder *isolate.IsolateCommandBuilder
	command               string
	submission            *model.Submission
}

func NewTestCaseExecutor(logger *zap.Logger, isolateDir string, command string, submission *model.Submission) *TestCaseExecutor {
	return &TestCaseExecutor{
		logger:     logger,
		isolateDir: isolateDir,
		command:    command,
		submission: submission,
		isolateCommandBuilder: isolate.NewIsolateCommandBuilder().
			WithProcesses(4).
			WithWallTime(submission.TimeLimit + 4).
			WithMaxFileSize(5120).
			AddDir("/etc:noexec").
			AddDir(isolateDir).
			WithCGroup().
			WithTime(submission.TimeLimit).
			WithExtraTime(submission.TimeLimit).
			WithCGroupMemory(submission.MemoryLimit).
			WithStackSize(submission.MemoryLimit).
			AddEnv("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").
			WithStderrToStdout(),
	}
}

func (e *TestCaseExecutor) Execute(testcase *model.TestCase) *model.SubmissionResult {

	boxId := snowflakeid.NewInt()%999 + 1
	if err := e.setupEnvironment(testcase, boxId); err != nil {
		return e.handleError(err, boxId)
	}
	defer e.cleanup(boxId, testcase)

	execResult, err := e.runIsolatedCommand(testcase, boxId)
	if err != nil {
		return e.handleError(err, boxId)
	}

	return e.processExecutionResult(execResult, testcase)
}

func (e *TestCaseExecutor) setupEnvironment(testcase *model.TestCase, boxId int) error {
	inputFilename := e.isolateDir + "/" + testcase.GetInputFileName()
	commandShFile := e.isolateDir + "/" + snowflakeid.NewString() + ".sh"

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
	commandShFile := e.isolateDir + "/" + snowflakeid.NewString() + ".sh"

	args := e.isolateCommandBuilder.Clone().
		WithBoxID(boxId).
		WithStdinFile(inputFilename).
		WithStdoutFile(outputFilename).
		WithMetaFile(metaOutFilename).
		WithRunCommands("/bin/bash", commandShFile).
		Build()

	execCmd := exec.Command(args[0], args[1:]...)
	execCmd.Dir = e.isolateDir

	if err := execCmd.Run(); err != nil {
		return nil, err
	}

	return isolate.NewMetaResultFromFile(metaOutFilename)
}

func (e *TestCaseExecutor) processExecutionResult(metaResult *isolate.MetaResult, testcase *model.TestCase) *model.SubmissionResult {
	result := &model.SubmissionResult{
		SubmissionID: e.submission.ID,
		TestCaseID:   testcase.ID,
		MemoryUsage:  metaResult.CGMem,
		TimeUsage:    metaResult.TimeWall,
	}

	if metaResult.ExitCode != 0 {
		result.Status = &StatusRuntimeError
		result.Stdout = &metaResult.Message
		return result
	}

	if metaResult.TimeWall > float64(e.submission.TimeLimit) {
		result.Status = &StatusTimeLimitExceeded
		return result
	}

	if metaResult.CGMem > float64(e.submission.MemoryLimit) {
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

	if output != *testcase.ExpectOutput {
		result.Status = &StatusWrongAnswer
		result.Stdout = &output
		return result
	}

	result.Status = &StatusSuccess
	result.Stdout = &output
	return result
}

func (e *TestCaseExecutor) readOutput(filename string) (string, error) {
	file, err := os.Open(e.isolateDir + "/" + filename)
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

func (e *TestCaseExecutor) handleError(err error, boxId int) *model.SubmissionResult {
	e.logger.Error("Execution error",
		zap.Int("boxId", boxId),
		zap.Error(err))

	errStr := err.Error()
	return &model.SubmissionResult{
		SubmissionID: e.submission.ID,
		Status:       &StatusRuntimeError,
		Stdout:       &errStr,
	}
}
