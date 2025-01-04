package executor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing/iotest"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/isolate"
	snowflakeid "github.com/Mirai3103/remote-compiler/pkg/snowflake_id"
	"go.uber.org/zap"
)

type Executor interface {
	Compile(*model.Submission) error
	Execute(*model.Submission, chan<- *model.SubmissionResult) error
}

type executor struct {
	logger *zap.Logger
	cfg    config.ExecutorConfig
}

var (
	SuccessStatus             = "Success"
	CompileErrorStatus        = "Compile Error"
	RuntimeErrorStatus        = "Runtime Error"
	WrongAnswerStatus         = "Wrong Answer"
	TimeLimitExceededStatus   = "Time Limit Exceeded"
	MemoryLimitExceededStatus = "Memory Limit Exceeded"
)

func (e *executor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	testcase := submission.TestCases
	command := strings.ReplaceAll(*submission.Language.RunCommand, "$BinaryFileName", e.cfg.IsolateDir+"/"+submission.Language.GetBinaryFileName())
	command = strings.ReplaceAll(command, "$SourceFileName", e.cfg.IsolateDir+"/"+submission.Language.GetSourceFileName())
	// defer os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetSourceFileName())
	// defer os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetBinaryFileName())
	wg := sync.WaitGroup{}

	isolateCommandBuilder := isolate.NewIsolateCommandBuilder().WithProcesses(4).WithWallTime(submission.TimeLimit + 4).WithMaxFileSize(5120).AddDir("/etc:noexec").AddDir(e.cfg.IsolateDir).WithCGroup().WithTime(submission.TimeLimit).WithExtraTime(submission.TimeLimit).WithCGroupMemory(submission.MemoryLimit).WithStackSize(submission.MemoryLimit).AddEnv("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").WithStderrToStdout()

	for _, testcase := range testcase {
		wg.Add(1)
		go func(testcase *model.TestCase) {
			defer wg.Done()
			result := &model.SubmissionResult{
				SubmissionID: submission.ID,
				TestCaseID:   testcase.ID,
			}
			inputFilename := e.cfg.IsolateDir + "/" + testcase.GetInputFileName()
			outputFilename := testcase.GetExpectOutputFileName()
			metaOutFilename := e.cfg.IsolateDir + "/" + testcase.GetExpectOutputFileName() + ".meta"
			os.WriteFile(inputFilename, []byte(*testcase.Input), 0644)
			commandShFile := e.cfg.IsolateDir + "/" + snowflakeid.NewString() + ".sh"
			os.WriteFile(commandShFile, []byte(command), 0644)
			// defer os.Remove(commandShFile)
			// defer os.Remove(inputFilename)
			// defer os.Remove(outputFilename)
			// defer os.Remove(metaOutFilename)
			boxId := snowflakeid.NewInt()
			boxDir, _ := isolate.InitBox(boxId)
			// defer func(boxId int) {
			// 	isolate.CleanBox(boxId)
			// }(boxId)

			args := isolateCommandBuilder.Clone().WithBoxID(boxId).WithStdinFile(inputFilename).WithStdoutFile(outputFilename).WithMetaFile(metaOutFilename).WithRunCommands("/bin/bash", commandShFile).Build()
			execCmd := exec.Command(args[0], args[1:]...)
			var stderr strings.Builder
			execCmd.Stderr = &stderr
			fmt.Println(*submission.Code)
			e.logger.Info("Execute command: ", zap.Any("args", args))
			execCmd.Dir = e.cfg.IsolateDir
			err := execCmd.Run()
			file, _ := os.Open(*boxDir + "/" + outputFilename)
			defer file.Close()

			b, _ := io.ReadAll(iotest.OneByteReader(file))
			if err != nil {
				errStr := err.Error()
				stdout := string(b) + "\n" + stderr.String()
				e.logger.Error("Error when run boxId: ", zap.Any("boxId", boxId), zap.Error(err), zap.String("stdout", errStr))
				result.Status = &RuntimeErrorStatus
				result.Stdout = &stdout
				ch <- result
				return
			}
			metaResult, err := isolate.NewMetaResultFromFile(metaOutFilename)
			if err != nil {
				errStr := err.Error()
				e.logger.Error("Error when read meta of boxId: ", zap.Any("boxId", boxId), zap.Error(err))
				result.Status = &RuntimeErrorStatus
				result.Stdout = &errStr
				ch <- result
				return
			}
			if metaResult.ExitCode != 0 {
				result.Status = &RuntimeErrorStatus
				result.Stdout = &metaResult.Message
				ch <- result
				return
			}
			result.MemoryUsage = metaResult.CGMem
			result.TimeUsage = metaResult.TimeWall

			if metaResult.TimeWall > float64(submission.TimeLimit) {
				result.Status = &TimeLimitExceededStatus
				ch <- result
				return
			}

			if metaResult.CGMem > float64(submission.MemoryLimit) {
				result.Status = &MemoryLimitExceededStatus
				ch <- result
				return
			}

			if err != nil {
				errStr := err.Error()
				e.logger.Error("Error when read output file of boxId: ", zap.Any("boxId", boxId), zap.Error(err))
				ch <- &model.SubmissionResult{
					SubmissionID: submission.ID,
					TestCaseID:   testcase.ID,
					Status:       &RuntimeErrorStatus,
					Stdout:       &errStr,
				}
				return
			}
			stdout := string(b)
			expectedOutput := *testcase.ExpectOutput
			if stdout != expectedOutput {
				result.Status = &WrongAnswerStatus
				result.Stdout = &stdout
				ch <- result
				return
			}
			result.Status = &SuccessStatus
			result.Stdout = &stdout
			e.logger.Debug("Result of boxId: ", zap.Any("boxId", boxId), zap.Any("result", result))
			ch <- result
		}(&testcase)
	}
	wg.Wait()
	close(ch)
	return nil
}

// Compile implements Executor.
func (e *executor) Compile(sb *model.Submission) error {
	language := sb.Language
	code := sb.Code
	sourceFilename := e.cfg.IsolateDir + "/" + language.GetSourceFileName()
	err := os.WriteFile(sourceFilename, []byte(*code), 0644)

	if err != nil {
		e.logger.Error("Error when write source file", zap.Error(err))
		return err
	}
	if sb.Language.CompileCommand == nil || *sb.Language.CompileCommand == "" {
		return nil
	}

	binaryFilename := e.cfg.IsolateDir + "/" + language.GetBinaryFileName()

	command := strings.ReplaceAll(*language.CompileCommand, "$SourceFileName", sourceFilename)
	command = strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)
	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.cfg.IsolateDir
	var stderr strings.Builder
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		e.logger.Error("Error when compile", zap.Error(err), zap.String("stderr", stderr.String()))
		return errors.New(stderr.String())
	}
	os.Remove(sourceFilename)
	return nil

}

func NewExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	ex := &executor{
		logger: logger,
		cfg:    cfg,
	}
	return ex
}
