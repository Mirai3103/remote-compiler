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

func (e *executor) initIsolate() {
	execCmd := exec.Command("isolate", "--cg", "--init")
	execCmd.Run()
}

var (
	SuccessStatus             = "Success"
	CompileErrorStatus        = "Compile Error"
	RuntimeErrorStatus        = "Runtime Error"
	WrongAnswerStatus         = "Wrong Answer"
	TimeLimitExceededStatus   = "Time Limit Exceeded"
	MemoryLimitExceededStatus = "Memory Limit Exceeded"
)

// isolate --cg -t $timeLimit -x 1 -w ($timeLimit + 4) -k 64000 -p 5 --cg-mem=128000 -f 5120 -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -d /etc:noexec --run -- $RunCommand
func (e *executor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	testcase := submission.TestCases
	command := strings.ReplaceAll(*submission.Language.RunCommand, "$BinaryFileName", e.cfg.IsolateDir+"/"+submission.Language.GetBinaryFileName())
	command = strings.ReplaceAll(command, "$SourceFileName", e.cfg.IsolateDir+"/"+submission.Language.GetSourceFileName())
	wg := sync.WaitGroup{}

	isolateCommandBuilder := isolate.NewIsolateCommandBuilder().WithWallTime(submission.TimeLimit + 4).WithMaxFileSize(5120).AddDir("/etc:noexec").AddDir(e.cfg.IsolateDir).WithCGroup().WithTime(submission.TimeLimit).WithExtraTime(submission.TimeLimit).WithCGroupMemory(submission.MemoryLimit).WithStackSize(submission.MemoryLimit).AddEnv("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").WithStderrToStdout()

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

			args := isolateCommandBuilder.WithStdinFile(inputFilename).WithStdoutFile(outputFilename).WithMetaFile(metaOutFilename).WithRunCommand(command).Build()
			execCmd := exec.Command(args[0], args[1:]...)
			execCmd.Dir = e.cfg.IsolateDir
			fmt.Println(args)
			err := execCmd.Run()
			if err != nil {
				errStr := err.Error()
				result.Status = &RuntimeErrorStatus
				result.Stdout = &errStr
				ch <- result
				return
			}
			metaResult, err := isolate.NewMetaResultFromFile(metaOutFilename)
			if err != nil {
				errStr := err.Error()
				result.Status = &RuntimeErrorStatus
				result.Stdout = &errStr
				ch <- result
			}
			if metaResult.ExitCode != 0 {
				result.Status = &RuntimeErrorStatus
				result.Stdout = &metaResult.Message
				ch <- result
				return
			}

			if metaResult.TimeWall > float64(submission.TimeLimit) {
				result.Status = &TimeLimitExceededStatus
				ch <- result
				return
			}

			if metaResult.CGMem > submission.MemoryLimit {
				result.Status = &MemoryLimitExceededStatus
				ch <- result
				return
			}

			file, err := os.Open("/var/local/lib/isolate/0/box/" + outputFilename)
			if err != nil {
				errStr := err.Error()
				result.Status = &RuntimeErrorStatus
				result.Stdout = &errStr
				ch <- result
				return
			}
			defer file.Close()

			b, err := io.ReadAll(iotest.OneByteReader(file))
			if err != nil {
				errStr := err.Error()
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
			ch <- result
		}(&testcase)
	}
	wg.Wait()
	close(ch)
	return nil
}

// Compile implements Executor.
func (e *executor) Compile(sb *model.Submission) error {
	if sb.Language.CompileCommand == nil {
		return nil
	}
	language := sb.Language
	code := sb.Code
	sourceFilename := e.cfg.IsolateDir + "/" + language.GetSourceFileName()
	binaryFilename := e.cfg.IsolateDir + "/" + language.GetBinaryFileName()
	err := os.WriteFile(sourceFilename, []byte(*code), 0644)
	if err != nil {
		return err
	}

	command := strings.ReplaceAll(*language.CompileCommand, "$SourceFileName", sourceFilename)
	command = strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)
	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.cfg.IsolateDir
	var stderr strings.Builder
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	defer os.Remove(sourceFilename)
	return nil

}

func NewExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	ex := &executor{
		logger: logger,
		cfg:    cfg,
	}
	ex.initIsolate()
	return ex
}
