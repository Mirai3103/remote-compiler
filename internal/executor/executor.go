package executor

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"go.uber.org/zap"
)

type Executor interface {
	Compile(*model.Submission) error
	Execute(*model.Submission, chan<- *model.SubmissionResult) error
}

type executor struct {
	logger      *zap.Logger
	cfx         config.ExecutorConfig
	isolatePath string
}

func (e *executor) initIsolate() {
	execCmd := exec.Command("isolate", "--init")
	execCmd.Run()
	e.isolatePath = "/var/local/lib/isolate/0"
}

func (e *executor) run(command string) (*string, error) {
	// rs := exec.Command(*command)
	// err := rs.Run()
	// if err != nil {
	// 	return err
	// }
	// return nil
	randomSleep := rand.IntN(20)
	time.Sleep(time.Duration(randomSleep) * time.Second)
	fmt.Print("Sleeping for ", randomSleep, " seconds\n")
	return nil, nil
}

// isolate --cg -t $timeLimit -x 1 -w ($timeLimit + 4) -k 64000 -p 5 --cg-mem=128000 -f 5120 -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -d /etc:noexec --run -- $RunCommand
// Execute implements Executor.
func (e *executor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	testcase := submission.TestCases
	command := strings.ReplaceAll(*submission.Language.RunCommand, "$BinaryFileName", e.isolatePath+"/"+submission.Language.GetBinaryFileName())
	command = strings.ReplaceAll(*submission.Language.RunCommand, "$SourceFileName", e.isolatePath+"/"+submission.Language.GetSourceFileName())
	var wg sync.WaitGroup
	for _, testcase := range testcase {
		wg.Add(1)
		go func(testcase *model.TestCase) {
			defer wg.Done()
			isolateCommand := fmt.Sprintf("isolate --cg -t %d -x 1 -w %d -k 64000 -p 5 --cg-mem=128000 -f 5120 -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" -d /etc:noexec --run -- %s", submission.TimeLimit, submission.TimeLimit+4, command)
			cmd := exec.Command("sh", "-c", isolateCommand)
			var stdout strings.Builder
			var stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			result := &model.SubmissionResult{
				SubmissionID: submission.ID,
				TestCaseID:   testcase.ID,
			}
			if err != nil {
				status := "Runtime Error"
				errStr := stderr.String()
				result.Stderr = &errStr
				result.Status = &status
			} else {
				status := "Accepted"
				stdoutStr := stdout.String()
				result.Stdout = &stdoutStr
				result.Status = &status
			}

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
	sourceFilename := e.isolatePath + "/" + language.GetSourceFileName()
	binaryFilename := e.isolatePath + "/" + language.GetBinaryFileName()
	err := os.WriteFile(sourceFilename, []byte(*code), 0644)
	if err != nil {
		return err
	}

	command := strings.ReplaceAll(*language.CompileCommand, "$SourceFileName", sourceFilename)
	command = strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)
	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.isolatePath
	var stderr strings.Builder
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	defer os.Remove(sourceFilename)
	return nil

}

func NewExecutor(logger *zap.Logger, cfx config.ExecutorConfig) Executor {
	ex := &executor{
		logger: logger,
		cfx:    cfx,
	}
	ex.initIsolate()
	return ex
}
